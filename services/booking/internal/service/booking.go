package service

import (
	"context"
	"fmt"
	"time"

	"github.com/DipeshC1707/BookYourShow/booking/internal/inventoryclient"
	"github.com/DipeshC1707/BookYourShow/booking/internal/db"
)

type BookingService struct {
	inventory inventoryclient.Client
	db *db.DB
	lockTTL   time.Duration
}

func NewBookingService(
	inventory inventoryclient.Client,
	db *db.DB,
	lockTTL time.Duration,
) *BookingService {
	return &BookingService{
		inventory: inventory,
		db: db,
		lockTTL:   lockTTL,
	}
}

func (s *BookingService) CreateBooking(
	ctx context.Context,
	eventID string,
	seatIDs []string,
	userID string,
	idempotencyKey string,
) (string, error) {

	if eventID == "" {
		return "", fmt.Errorf("event id is required")
	}
	if len(seatIDs) == 0 {
		return "", fmt.Errorf("at least one seat is required")
	}
	if userID == "" {
		return "", fmt.Errorf("user id is required")
	}

	var existingBookingID string

	err := s.db.Pool.QueryRow(
		ctx,
		`
		SELECT booking_id
		FROM bookings
		WHERE idempotency_key = $1
		`,
		idempotencyKey,
	).Scan(&existingBookingID)

	if err == nil {
		// Booking already exists → idempotent return
		return existingBookingID, nil
	}

	// 1️⃣ Lock seats in Inventory
	if err := s.inventory.LockSeats(ctx, eventID, seatIDs, userID); err != nil {
		return "", err
	}

	bookingID := fmt.Sprintf("bkg_%d", time.Now().UnixNano())
	expiresAt := time.Now().Add(s.lockTTL)

	// 2️⃣ Start DB transaction
	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	// 3️⃣ Insert booking
	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO bookings (
			booking_id,
			user_id,
			event_id,
			seat_ids,
			status,
			expires_at,
			idempotency_key
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		`,
		bookingID,
		userID,
		eventID,
		seatIDs,
		"CREATED",
		expiresAt,
		idempotencyKey,
	)
	

	// 4️⃣ Insert booking seats
	for _, seatID := range seatIDs {
		_, err = tx.Exec(
			ctx,
			`
			INSERT INTO booking_seats (
				booking_id,
				event_id,
				seat_id,
				status
			)
			VALUES ($1, $2, $3, $4)
			`,
			bookingID,
			eventID,
			seatID,
			"CREATED",
		)
		if err != nil {
			return "", err
		}
	}

	// 5️⃣ Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return "", err
	}

	return bookingID, nil
}

func (s *BookingService) ConfirmBooking(
	ctx context.Context,
	bookingID string,
) error {

	var status string

	err := s.db.Pool.QueryRow(
		ctx,
		"SELECT status FROM bookings WHERE booking_id = $1",
		bookingID,
	).Scan(&status)

	if err != nil {
		return err
	}

	if status == "CONFIRMED" {
		return nil // idempotent success
	}

	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	cmd, err := tx.Exec(
		ctx,
		`
		UPDATE bookings
		SET status = $1, updated_at = now()
		WHERE booking_id = $2
		`,
		"CONFIRMED",
		bookingID,
	)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("booking not found")
	}

	_, err = tx.Exec(
		ctx,
		`
		UPDATE booking_seats
		SET status = $1
		WHERE booking_id = $2
		`,
		"CONFIRMED",
		bookingID,
	)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	// Redis cleanup is best-effort
	// TTL already protects us
	return nil
}

type ExpiredBooking struct {
	BookingID string
	EventID   string
	SeatIDs   []string
}

func (s *BookingService) FindExpiredBookings(
	ctx context.Context,
	limit int,
) ([]ExpiredBooking, error) {

	rows, err := s.db.Pool.Query(ctx, `
		SELECT b.booking_id, b.event_id, bs.seat_id
		FROM bookings b
		JOIN booking_seats bs ON b.booking_id = bs.booking_id
		WHERE b.status = 'CREATED'
		  AND b.expires_at < now()
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	m := map[string]*ExpiredBooking{}

	for rows.Next() {
		var bookingID, eventID, seatID string
		if err := rows.Scan(&bookingID, &eventID, &seatID); err != nil {
			return nil, err
		}

		if _, ok := m[bookingID]; !ok {
			m[bookingID] = &ExpiredBooking{
				BookingID: bookingID,
				EventID:   eventID,
			}
		}
		m[bookingID].SeatIDs = append(m[bookingID].SeatIDs, seatID)
	}

	result := make([]ExpiredBooking, 0, len(m))
	for _, b := range m {
		result = append(result, *b)
	}

	return result, nil
}


func (s *BookingService) CancelExpiredBooking(
	ctx context.Context,
	bookingID string,
) error {

	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	ct, err := tx.Exec(ctx, `
		UPDATE bookings
		SET status = 'CANCELLED'
		WHERE booking_id = $1
		  AND status = 'CREATED'
	`, bookingID)
	if err != nil {
		return err
	}

	// If nothing updated, booking was already handled
	if ct.RowsAffected() == 0 {
		return nil
	}

	_, err = tx.Exec(ctx, `
		UPDATE booking_seats
		SET status = 'CANCELLED'
		WHERE booking_id = $1
	`, bookingID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

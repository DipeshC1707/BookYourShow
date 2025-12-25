package service

import (
	"context"
	"fmt"
	"time"

	"github.com/DipeshC1707/BookYourShow/booking/internal/inventoryclient"
)

type BookingService struct {
	inventory inventoryclient.Client
	lockTTL   time.Duration
}

func NewBookingService(
	inventory inventoryclient.Client,
	lockTTL time.Duration,
) *BookingService {
	return &BookingService{
		inventory: inventory,
		lockTTL:   lockTTL,
	}
}

func (s *BookingService) CreateBooking(
	ctx context.Context,
	eventID string,
	seatIDs []string,
	userID string,
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

	if err := s.inventory.LockSeats(ctx, eventID, seatIDs, userID); err != nil {
		return "", err
	}

	bookingID := fmt.Sprintf("bkg_%d", time.Now().UnixNano())
	expiresAt := time.Now().Add(s.lockTTL)

	_ = expiresAt // used later when DB is added

	// TODO: save booking to DB

	return bookingID, nil
}

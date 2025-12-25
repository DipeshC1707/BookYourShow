package expiration

import (
	"context"
)

func (r *Runner) runOnce(ctx context.Context) {
	bookings, err := r.bookingService.FindExpiredBookings(ctx, 50)
	if err != nil {
		return
	}

	for _, b := range bookings {
		if err := r.bookingService.CancelExpiredBooking(ctx, b.BookingID); err != nil {
			continue
		}

		_ = r.inventory.ReleaseSeats(ctx, b.EventID, b.SeatIDs)
	}
}


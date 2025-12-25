package expiration

import (
	"context"
	"time"

	"github.com/DipeshC1707/BookYourShow/booking/internal/service"
	"github.com/DipeshC1707/BookYourShow/booking/internal/inventoryclient"
)

type Runner struct {
	bookingService *service.BookingService
	inventory      inventoryclient.Client
	interval       time.Duration
}

func NewRunner(
	bookingService *service.BookingService,
	inventory inventoryclient.Client,
	interval time.Duration,
) *Runner {
	return &Runner{
		bookingService: bookingService,
		inventory:      inventory,
		interval:       interval,
	}
}

func (r *Runner) Start(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.runOnce(ctx)
		case <-ctx.Done():
			return
		}
	}
}

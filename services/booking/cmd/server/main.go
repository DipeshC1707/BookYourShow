package main

import (
	"context"
	"time"

	"os"
	"github.com/DipeshC1707/BookYourShow/booking/internal/db"
	"github.com/DipeshC1707/BookYourShow/booking/internal/inventoryclient"
	"github.com/DipeshC1707/BookYourShow/booking/internal/service"
)

func main() {
	ctx := context.Background()

	inventory, err := inventoryclient.New("localhost:50051")
	if err != nil {
		panic(err)
	}
	defer inventory.Close()

	dbConn, err := db.New(ctx, os.Getenv("BOOKING_DB_DSN"))
	if err != nil {
	panic(err)
	}
	defer dbConn.Close()


	bookingService := service.NewBookingService(
		inventory,
		10*time.Minute,
	)

	// TEMP test
	_, err = bookingService.CreateBooking(
		ctx,
		"event1",
		[]string{"A1", "A2"},
		"user42",
	)
	if err != nil {
		panic(err)
	}

	select {} // keep process alive
}

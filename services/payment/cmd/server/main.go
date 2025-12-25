package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/DipeshC1707/BookYourShow/payment/internal/bookingclient"
)

func main() {
	bookingAddr := os.Getenv("BOOKING_GRPC_ADDR")
	if bookingAddr == "" {
		bookingAddr = "localhost:50052"
	}

	client, err := bookingclient.New(bookingAddr)
	if err != nil {
		log.Fatalf("failed to connect to booking: %v", err)
	}
	defer client.Close()

	// MOCK PAYMENT FLOW
	bookingID := os.Getenv("BOOKING_ID")
	if bookingID == "" {
		log.Fatal("BOOKING_ID is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("processing payment for booking %s", bookingID)

	// Simulate payment delay
	time.Sleep(2 * time.Second)

	if err := client.ConfirmBooking(ctx, bookingID); err != nil {
		log.Fatalf("payment failed: %v", err)
	}

	log.Printf("payment successful, booking confirmed: %s", bookingID)
}
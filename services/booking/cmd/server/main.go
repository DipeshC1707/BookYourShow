package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	bookingpb "github.com/DipeshC1707/BookYourShow/proto/booking/v1"
	"github.com/DipeshC1707/BookYourShow/booking/internal/db"
	"github.com/DipeshC1707/BookYourShow/booking/internal/grpcserver"
	"github.com/DipeshC1707/BookYourShow/booking/internal/inventoryclient"
	"github.com/DipeshC1707/BookYourShow/booking/internal/service"
	"github.com/DipeshC1707/BookYourShow/booking/internal/expiration"
)

func main() {
	// Root context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// -------------------------
	// Handle graceful shutdown
	// -------------------------
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// -------------------------
	// Load configuration
	// -------------------------
	inventoryAddr := os.Getenv("INVENTORY_GRPC_ADDR")
	if inventoryAddr == "" {
		inventoryAddr = "localhost:50051"
	}

	dbDSN := os.Getenv("BOOKING_DB_DSN")
	if dbDSN == "" {
		log.Fatal("BOOKING_DB_DSN is required")
	}

	grpcPort := os.Getenv("BOOKING_GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50052"
	}

	// -------------------------
	// Connect to Inventory
	// -------------------------
	inventoryClient, err := inventoryclient.New(inventoryAddr)
	if err != nil {
		log.Fatalf("failed to connect to inventory: %v", err)
	}
	defer inventoryClient.Close()

	// -------------------------
	// Connect to Database
	// -------------------------
	dbConn, err := db.New(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	  
	// -------------------------
	// Create Booking Service
	// -------------------------
	bookingService := service.NewBookingService(
		inventoryClient,
		dbConn,
		2*time.Minute, // seat lock TTL
	)

	expirer := expiration.NewRunner(
		bookingService,
		inventoryClient,
		30*time.Second,
	)
	
	go expirer.Start(ctx)
	
	// -------------------------
	// Start gRPC server
	// -------------------------
	grpcServer := grpc.NewServer()

	bookingGrpcServer := grpcserver.New(bookingService)
	bookingpb.RegisterBookingServiceServer(grpcServer, bookingGrpcServer)

	listener, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", grpcPort, err)
	}

	go func() {
		log.Printf("Booking gRPC server started on :%s", grpcPort)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	// -------------------------
	// Wait for shutdown signal
	// -------------------------
	<-sigCh
	log.Println("shutdown signal received")

	grpcServer.GracefulStop()
	cancel()

	log.Println("booking service stopped gracefully")
}

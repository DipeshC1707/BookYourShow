package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"net"

	"github.com/DipeshC1707/BookYourShow/inventory/internal/config"
	"github.com/DipeshC1707/BookYourShow/inventory/internal/logger"
	"github.com/DipeshC1707/BookYourShow/inventory/internal/server"
	"github.com/DipeshC1707/BookYourShow/inventory/internal/redis"
	"os"

	"google.golang.org/grpc"

	"github.com/DipeshC1707/BookYourShow/inventory/internal/grpcserver"
	grpcpb "github.com/DipeshC1707/BookYourShow/proto/inventory/v1"
	"github.com/DipeshC1707/BookYourShow/inventory/internal/service"
)

func main() {
	// 1️⃣ Load configuration
	cfg := config.Load()

	// 2️⃣ Create logger
	log := logger.New(cfg.ServiceName)

	// 3️⃣ Root context (cancelled on shutdown)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 4️⃣ Connect to Redis (fail fast)
	redisClient := redisclient.New(cfg.RedisURL)
	if err := redisClient.Ping(ctx); err != nil {
		log.Error("failed to connect to redis", "error", err)
		os.Exit(1)
	}
	log.Info("connected to redis")

	// 5️⃣ Create Inventory business service
	inventoryService := service.NewInventoryService(
		redisClient,
		10*time.Minute, // seat lock TTL
	)

	// 6️⃣ Start gRPC server
	grpcSrv := grpc.NewServer()

	inventoryGrpcServer := grpcserver.NewServer(inventoryService)
	grpcpb.RegisterInventoryServiceServer(
		grpcSrv,
		inventoryGrpcServer,
	)

	grpcListener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Error("failed to listen on gRPC port", "error", err)
		os.Exit(1)
	}

	go func() {
		log.Info("inventory gRPC server started", "port", 50051)
		if err := grpcSrv.Serve(grpcListener); err != nil {
			log.Error("gRPC server error", "error", err)
		}
	}()

	// 7️⃣ HTTP server (health checks only)
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Info("health server started", "port", cfg.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("HTTP server error", "error", err)
		}
	}()

	// 8️⃣ Wait for shutdown signal (SIGTERM / Ctrl+C)
	server.WaitForShutdown(
		time.Duration(cfg.ShutdownTimeoutSeconds)*time.Second,
		cancel,
	)

	log.Info("shutdown initiated")

	// 9️⃣ Graceful shutdown
	grpcSrv.GracefulStop()

	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		time.Duration(cfg.ShutdownTimeoutSeconds)*time.Second,
	)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error("HTTP shutdown failed", "error", err)
	}

	if err := redisClient.Close(); err != nil {
		log.Error("failed to close redis", "error", err)
	}

	log.Info("inventory service shutdown complete")
}
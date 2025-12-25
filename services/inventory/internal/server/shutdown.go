package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func WaitForShutdown(timeout time.Duration, cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	ctx, shutdownCancel := context.WithTimeout(context.Background(), timeout)
	defer shutdownCancel()

	cancel()

	<-ctx.Done()
}

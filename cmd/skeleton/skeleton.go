package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kianooshaz/skeleton/internal/container"
)

func main() {
	if err := run(); err != nil {
		slog.Error("Application failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	// Create cancellable context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize web container
	c, err := container.NewWebContainer()
	if err != nil {
		return fmt.Errorf("initializing container: %w", err)
	}

	// Start all services
	if err := c.Start(cancel); err != nil {
		cancel()
	}

	c.Logger().Info("Application started successfully")

	waitForShutdown(ctx, c)

	c.Logger().Info("Application shut down gracefully")

	return nil
}

func waitForShutdown(ctx context.Context, c container.Container) {
	// Setup signal handling
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer signal.Stop(shutdown)

	// Wait for shutdown signal
	select {
	case sig := <-shutdown:
		c.Logger().Info("Received shutdown signal", "signal", sig)
	case <-ctx.Done():
		c.Logger().Info("Context cancelled, shutting down")
	}

	if err := c.Stop(); err != nil {
		c.Logger().Error("Error stopping container", "error", err)
	}
}

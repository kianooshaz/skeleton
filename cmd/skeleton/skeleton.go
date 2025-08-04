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

	// Ensure proper cleanup
	defer func() {
		if closeErr := c.Close(); closeErr != nil {
			slog.Error("Failed to close container", "error", closeErr)
		}
	}()

	// Start all services
	if err := c.Start(); err != nil {
		return fmt.Errorf("starting services: %w", err)
	}

	c.Logger.Info("Application started successfully")

	// Wait for shutdown signal
	if err := waitForShutdown(ctx, c); err != nil {
		return fmt.Errorf("shutdown failed: %w", err)
	}

	c.Logger.Info("Application shut down successfully")
	return nil
}

func waitForShutdown(ctx context.Context, c *container.WebContainer) error {
	// Setup signal handling
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer signal.Stop(shutdown)

	// Wait for shutdown signal
	select {
	case sig := <-shutdown:
		c.Logger.Info("Received shutdown signal", "signal", sig.String())
	case <-ctx.Done():
		c.Logger.Info("Context cancelled, shutting down")
	}

	// Create timeout context for graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), c.Config.ShutdownTimeout)
	defer cancel()

	// Perform graceful shutdown
	if err := c.Stop(shutdownCtx); err != nil {
		return fmt.Errorf("stopping container: %w", err)
	}

	return nil
}

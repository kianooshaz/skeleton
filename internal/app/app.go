package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kianooshaz/skeleton/internal/container"

	_ "github.com/kianooshaz/skeleton/foundation/database/postgres"
	_ "github.com/kianooshaz/skeleton/foundation/log"
)

func Run() error {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Initialize web dependency injection container
	c, err := container.NewWebContainer()
	if err != nil {
		return fmt.Errorf("initializing web container: %w", err)
	}
	defer func() {
		if closeErr := c.Close(); closeErr != nil {
			c.Logger.Error("Error closing web container", "error", closeErr)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	errServer := make(chan error, 1)
	go func() {
		errServer <- c.RestServer.Start()
	}()

	select {
	case err = <-errServer:
		return err
	case <-ctx.Done():
		c.Logger.Info("shutting down server")

		ctxWithTimeout, cancel := context.WithTimeout(context.Background(), c.Config.ShutdownTimeout)
		defer cancel()

		if err := c.RestServer.Shutdown(ctxWithTimeout); err != nil {
			errs := errors.Join(fmt.Errorf("shutting down server: %w", err))

			if err := c.RestServer.Close(); err != nil {
				errs = errors.Join(errs, fmt.Errorf("closing server: %w", err))
			}

			return errs
		}

		stop()
	}

	c.Logger.Info("shutting down successfully")

	return nil
}

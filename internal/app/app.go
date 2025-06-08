package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kianooshaz/skeleton/foundation/config"
	"github.com/kianooshaz/skeleton/foundation/database/postgres"
	"github.com/kianooshaz/skeleton/internal/app/web/rest"

	_ "github.com/kianooshaz/skeleton/foundation/database/postgres"
	_ "github.com/kianooshaz/skeleton/foundation/log"
)

type Config struct {
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" validate:"required"`
	RestServer      rest.Config   `yaml:"rest_server" validate:"required"`
}

func Run() error {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg, err := config.Load[Config]("app")
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	restServer, err := rest.New(cfg.RestServer)
	if err != nil {
		return fmt.Errorf("creating server: %w", err)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	errServer := make(chan error, 1)
	go func() {
		errServer <- restServer.Start()
	}()

	select {
	case err = <-errServer:
		return err
	case <-ctx.Done():
		fmt.Println("shutting down server")

		ctxWithTimeout, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		if err := restServer.Shutdown(ctxWithTimeout); err != nil {
			errs := errors.Join(fmt.Errorf("shutting down server: %w", err))

			if err := restServer.Close(); err != nil {
				errs = errors.Join(errs, fmt.Errorf("closing server: %w", err))
			}

			return errs
		}

		if err := postgres.ConnectionPool.Close(); err != nil {
			return fmt.Errorf("closing postgres pool: %w", err)
		}

		stop()
	}

	fmt.Println("shutting down successfully")

	return nil
}

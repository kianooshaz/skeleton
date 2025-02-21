package web

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kianooshaz/skeleton/foundation/config"
	"github.com/kianooshaz/skeleton/foundation/log"
	"github.com/kianooshaz/skeleton/foundation/storage/postgres"
	"github.com/kianooshaz/skeleton/internal/app/web/rest"
)

type Config struct {
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" validate:"required"`
}

func Serve(configPath string) error {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	config.Init(configPath)

	cfg, err := config.Load[Config]("app.web")
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	log.Init()

	if err := postgres.Init(ctx); err != nil {
		return fmt.Errorf("creating postgres pool: %w", err)
	}

	if err := rest.Init(); err != nil {
		return fmt.Errorf("creating server: %w", err)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	errServer := make(chan error, 1)
	go func() {
		errServer <- rest.Server.Start()
	}()

	select {
	case err = <-errServer:
		return err
	case <-ctx.Done():
		fmt.Println("shutting down server")

		ctxWithTimeout, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		if err := rest.Server.Shutdown(ctxWithTimeout); err != nil {
			errs := errors.Join(fmt.Errorf("shutting down server: %w", err))

			if err := rest.Server.Close(); err != nil {
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

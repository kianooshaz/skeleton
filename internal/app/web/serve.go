package web

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kianooshaz/skeleton/foundation/database/postgres"
	"github.com/kianooshaz/skeleton/internal/app/web/rest"
	"github.com/kianooshaz/skeleton/internal/app/web/rest/handler"
	usernameService "github.com/kianooshaz/skeleton/modules/user/username/service"
)

func Serve(configPath string) error {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg, err := newConfig(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	db, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		return fmt.Errorf("creating postgres pool: %w", err)
	}

	_ = usernameService.New(cfg.UsernameService, db)

	server := rest.New(cfg.Rest, &handler.Handler{})

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	errServer := make(chan error, 1)
	go func() {
		errServer <- server.Start()
	}()

	select {
	case err = <-errServer:
		return err
	case <-ctx.Done():
		fmt.Println("shutting down server")

		ctxWithTimeout, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctxWithTimeout); err != nil {
			fmt.Println("failed to shutdown server", "error", err)
		}

		db.Close()

		stop()
	}

	fmt.Println("shutting down successfully")

	return nil
}

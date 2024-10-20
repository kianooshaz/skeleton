package web

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kianooshaz/skeleton/foundation/postgres"
	"github.com/kianooshaz/skeleton/internal/app/web/rest"
	"github.com/kianooshaz/skeleton/internal/app/web/rest/handler"
	"github.com/kianooshaz/skeleton/service/usernamesrv"
	"github.com/kianooshaz/skeleton/service/usersrv"
)

func Serve(configPath string) error {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg, err := newConfig(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	pgPool, err := postgres.NewPool(ctx, cfg.Postgres)
	if err != nil {
		return fmt.Errorf("creating postgres pool: %w", err)
	}

	userService := usersrv.New(pgPool)
	usernameService := usernamesrv.New(cfg.UsernameService, pgPool)

	server := rest.New(cfg.Rest, &handler.Handler{
		UserService:     userService,
		UsernameService: usernameService,
	})

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

		pgPool.Close()

		stop()
	}

	fmt.Println("shutting down successfully")

	return nil
}

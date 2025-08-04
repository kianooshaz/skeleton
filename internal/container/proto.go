package container

import (
	"context"
	"log/slog"
)

type Container interface {
	Start(cancel context.CancelFunc) error
	Stop() error
	Logger() *slog.Logger
}

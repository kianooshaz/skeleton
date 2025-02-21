package protocol

import "context"

type WebService interface {
	Start() error
	Shutdown(ctx context.Context) error
	Close() error
}

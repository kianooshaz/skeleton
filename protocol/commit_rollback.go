package protocol

import "context"

type CommitRollbacker interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

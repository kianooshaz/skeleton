package service

import (
	"context"
	"log/slog"
	"sync"

	"github.com/kianooshaz/skeleton/foundation/config"
	"github.com/kianooshaz/skeleton/foundation/database/postgres"
	"github.com/kianooshaz/skeleton/services/risk/audit/protocol"
	ap "github.com/kianooshaz/skeleton/services/risk/audit/protocol"
	"github.com/kianooshaz/skeleton/services/risk/audit/service/storage"
)

type (
	Config struct {
		BufferSize  int
		WorkerCount int
	}

	auditService struct {
		storage  storer
		logger   *slog.Logger
		recordCh chan ap.Record
		shutdown chan struct{}
		workerWg *sync.WaitGroup
	}

	storer interface {
		Create(ctx context.Context, record protocol.Record) error
	}
)

var Audit ap.Audit

func init() {
	cfg, err := config.Load[Config]("risk.audit")
	if err != nil {
		panic(err)
	}

	service := &auditService{
		storage: &storage.AuditStorage{Conn: postgres.ConnectionPool},
		logger: slog.With(
			slog.Group("package_info",
				slog.String("module", "user"),
				slog.String("service", "user"),
			),
		),
		recordCh: make(chan ap.Record, cfg.BufferSize),
		shutdown: make(chan struct{}),
		workerWg: &sync.WaitGroup{},
	}

	for range cfg.WorkerCount {

		go service.processRecords()
	}

	Audit = service
}

package service

import (
	"context"
	"database/sql"
	"log/slog"
	"sync"

	"github.com/kianooshaz/skeleton/foundation/config"
	"github.com/kianooshaz/skeleton/foundation/database/postgres"
	"github.com/kianooshaz/skeleton/foundation/order"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	"github.com/kianooshaz/skeleton/services/risk/audit/persistence"
	auditproto "github.com/kianooshaz/skeleton/services/risk/audit/proto"
)

type (
	Config struct {
		BufferSize  int
		WorkerCount int
	}

	persister interface {
		Create(ctx context.Context, record auditproto.Record) error
		Get(ctx context.Context, id auditproto.RecordID) (auditproto.Record, error)
		List(ctx context.Context, page pagination.Page, orderBy order.OrderBy) ([]auditproto.Record, error)
		Count(ctx context.Context) (int, error)
	}

	auditService struct {
		persister persister
		logger    *slog.Logger
		recordCh  chan auditproto.Record
		shutdown  chan struct{}
		workerWg  *sync.WaitGroup
		dbConn    *sql.DB
	}
)

var Service auditproto.AuditService

func init() {
	cfg, err := config.Load[Config]("risk.audit")
	if err != nil {
		panic(err)
	}

	service := &auditService{
		persister: &persistence.AuditStorage{Conn: postgres.ConnectionPool},
		logger: slog.With(
			slog.Group("package_info",
				slog.String("module", "audit"),
				slog.String("service", "audit"),
			),
		),
		recordCh: make(chan auditproto.Record, cfg.BufferSize),
		shutdown: make(chan struct{}),
		workerWg: &sync.WaitGroup{},
		dbConn:   postgres.ConnectionPool,
	}

	for range cfg.WorkerCount {
		go service.processRecords()
	}

	Service = service
}

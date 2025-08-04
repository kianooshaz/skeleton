package auditservice

import (
	"context"
	"database/sql"
	"log/slog"
	"sync"

	"github.com/kianooshaz/skeleton/foundation/order"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	"github.com/kianooshaz/skeleton/services/risk/audit/persistence"
	auditproto "github.com/kianooshaz/skeleton/services/risk/audit/proto"
)

type (
	Config struct {
		BufferSize  int `yaml:"buffer_size"`
		WorkerCount int `yaml:"worker_count"`
	}

	persister interface {
		Create(ctx context.Context, record auditproto.Record) error
		Get(ctx context.Context, id auditproto.RecordID) (auditproto.Record, error)
		List(ctx context.Context, page pagination.Page, orderBy order.OrderBy) ([]auditproto.Record, error)
		Count(ctx context.Context) (int, error)
	}

	Service struct {
		config    Config
		persister persister
		logger    *slog.Logger
		recordCh  chan auditproto.Record
		shutdown  chan struct{}
		workerWg  *sync.WaitGroup
		dbConn    *sql.DB
	}
)

// New creates a new audit service instance.
func New(cfg Config, db *sql.DB, logger *slog.Logger) auditproto.AuditService {
	serviceLogger := logger.With(
		slog.Group("package_info",
			slog.String("module", "audit"),
			slog.String("service", "audit"),
		),
	)

	// Set default values if not configured
	if cfg.BufferSize == 0 {
		cfg.BufferSize = 1000
	}
	if cfg.WorkerCount == 0 {
		cfg.WorkerCount = 3
	}

	svc := &Service{
		config:    cfg,
		persister: &persistence.AuditStorage{Conn: db},
		logger:    serviceLogger,
		recordCh:  make(chan auditproto.Record, cfg.BufferSize),
		shutdown:  make(chan struct{}),
		workerWg:  &sync.WaitGroup{},
		dbConn:    db,
	}

	// Start worker goroutines
	for range cfg.WorkerCount {
		go svc.processRecords()
	}

	return svc
}

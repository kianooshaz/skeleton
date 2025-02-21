// Package asyncaudit provides an audit service that records audit logs into a PostgreSQL database.
//
// Database Requirements:
// - The database must be PostgreSQL.
// - It should contain a table named `audit_records` with the following columns:
//   - id (UUID or appropriate type for unique identifiers)
//   - action (TEXT or VARCHAR)
//   - created_at (TIMESTAMP)
//   - data (JSONB or TEXT depending on the use case)
//   - origin_ip (TEXT or VARCHAR)
//   - resource_id (TEXT or VARCHAR)
//   - resource_type (TEXT or VARCHAR)
//   - user_id (TEXT or UUID depending on the user identifier type)
//
// The service ensures a graceful shutdown by processing all buffered audit records before termination.

package asyncaudit

import (
	"context"
	"database/sql"
	"sync"

	"github.com/kianooshaz/skeleton/foundation/audit"
	"github.com/kianooshaz/skeleton/foundation/log"
)

type AuditService struct {
	db       *sql.DB
	logger   log.Protocol
	recordCh chan audit.Record
	shutdown chan struct{}
	workerWg *sync.WaitGroup
}

type Config struct {
	BufferSize  int
	WorkerCount int
}

func NewAuditService(db *sql.DB, logger log.Protocol, cfg Config) *AuditService {
	service := &AuditService{
		db:       db,
		logger:   logger,
		recordCh: make(chan audit.Record, cfg.BufferSize),
		shutdown: make(chan struct{}),
		workerWg: &sync.WaitGroup{},
	}

	for range cfg.WorkerCount {

		go service.processRecords()
	}

	return service
}

func (as *AuditService) Record(record audit.Record) {
	as.workerWg.Add(1)
	as.recordCh <- record
}

func (as *AuditService) processRecords() {
	defer as.workerWg.Done()

	for {
		select {
		case record := <-as.recordCh:
			as.insertRecord(record)
		case <-as.shutdown:
			as.logger.Info("worker shutting down")
			return
		}
	}
}

func (as *AuditService) insertRecord(record audit.Record) {
	query := `INSERT INTO audit_records (id, action, created_at, data, origin_ip, resource_id, resource_type, user_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := as.db.Exec(query, record.ID, record.Action, record.CreatedAt, record.Data,
		record.OriginIP, record.ResourceID, record.ResourceType, record.UserID)
	if err != nil {
		as.logger.Error("failed to insert audit record", "error", err, "record", record)
	}
}

func (as *AuditService) Shutdown(ctx context.Context) {
	close(as.shutdown)

	done := make(chan struct{})
	go func() {
		as.workerWg.Wait()
		close(done)
	}()

	select {
	case <-done:
		as.logger.Info("all workers shut down gracefully")
	case <-ctx.Done():
		as.logger.Error("shutdown timeout; workers did not finish in time")
	}
}

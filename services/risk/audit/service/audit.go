package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	"github.com/kianooshaz/skeleton/foundation/session"
	auditproto "github.com/kianooshaz/skeleton/services/risk/audit/proto"
)

func (as *auditService) Record(record auditproto.Record) {
	// Generate ID if not provided
	if record.ID == auditproto.RecordID(uuid.Nil) {
		record.ID = auditproto.RecordID(uuid.New())
	}

	as.workerWg.Add(1)
	as.recordCh <- record
}

func (as *auditService) Get(ctx context.Context, req auditproto.GetRequest) (auditproto.GetResponse, error) {
	record, err := as.persister.Get(ctx, req.ID)
	if err != nil {
		return auditproto.GetResponse{}, err
	}

	return auditproto.GetResponse{Data: record}, nil
}

func (as *auditService) List(ctx context.Context, req auditproto.ListRequest) (auditproto.ListResponse, error) {
	records, err := as.persister.List(ctx, req.Page, req.OrderBy)
	if err != nil {
		return auditproto.ListResponse{}, err
	}

	count, err := as.persister.Count(ctx)
	if err != nil {
		return auditproto.ListResponse{}, err
	}

	return auditproto.ListResponse{
		Response: pagination.NewResponse(req.Page, count, records),
	}, nil
}

func (as *auditService) processRecords() {
	defer as.workerWg.Done()

	for {
		select {
		case record := <-as.recordCh:
			if err := as.persister.Create(context.Background(), record); err != nil {
				as.logger.ErrorContext(
					session.SetRequestID(context.Background(), record.RequestID),
					"failed to create audit record",
					slog.String("error", err.Error()),
					slog.Any("record", record),
				)
			}
		case <-as.shutdown:
			as.logger.Info("worker shutting down")
			return
		}
	}
}

func (as *auditService) Shutdown(ctx context.Context) {
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

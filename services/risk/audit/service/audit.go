package service

import (
	"context"
	"log/slog"

	"github.com/kianooshaz/skeleton/foundation/session"
	ap "github.com/kianooshaz/skeleton/services/risk/audit/protocol"
)

func (as *auditService) Record(record ap.Record) {
	as.workerWg.Add(1)
	as.recordCh <- record
}

func (as *auditService) processRecords() {
	defer as.workerWg.Done()

	for {
		select {
		case record := <-as.recordCh:
			if err := as.storage.Create(context.Background(), record); err != nil {
				as.logger.ErrorContext(session.SetRequestID(context.Background(), record.RequestID), "failed to create audit record", slog.String("error", err.Error()), "record", record)
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

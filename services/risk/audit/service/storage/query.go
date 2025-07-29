package storage

import (
	"context"

	dbproto "github.com/kianooshaz/skeleton/foundation/database/proto"
	"github.com/kianooshaz/skeleton/foundation/session"
	ap "github.com/kianooshaz/skeleton/services/risk/audit/protocol"
)

type AuditStorage struct {
	Conn dbproto.QueryExecutor
}

const create = `
	INSERT INTO 
		audit_records (request_id, action, created_at, data, origin_ip, resource_id, resource_type, user_id)
	VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8)
`

func (as *AuditStorage) Create(ctx context.Context, record ap.Record) error {
	conn := session.GetDBConnection(ctx, as.Conn)

	_, err := conn.ExecContext(ctx, create, record.RequestID, record.Action, record.CreatedAt, record.Data,
		record.OriginIP, record.ResourceID, record.ResourceType, record.UserID)
	return err
}

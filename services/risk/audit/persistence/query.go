package persistence

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"

	dbproto "github.com/kianooshaz/skeleton/foundation/database/proto"
	"github.com/kianooshaz/skeleton/foundation/derror"
	"github.com/kianooshaz/skeleton/foundation/order"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	"github.com/kianooshaz/skeleton/foundation/session"
	auditproto "github.com/kianooshaz/skeleton/services/risk/audit/proto"
)

type AuditStorage struct {
	Conn dbproto.QueryExecutor
}

const defaultPageSize = 20

//go:embed queries/create.sql
var createQuery string

//go:embed queries/get.sql
var getQuery string

//go:embed queries/list.sql
var listQuery string

//go:embed queries/count.sql
var countQuery string

func (as *AuditStorage) Create(ctx context.Context, record auditproto.Record) error {
	conn := session.GetDBConnection(ctx, as.Conn)

	_, err := conn.ExecContext(ctx, createQuery, record.ID, record.RequestID, record.Action,
		record.CreatedAt, record.Data, record.OriginIP, record.ResourceID, record.ResourceType, record.UserID)
	return err
}

func (as *AuditStorage) Get(ctx context.Context, id auditproto.RecordID) (auditproto.Record, error) {
	conn := session.GetDBConnection(ctx, as.Conn)

	row := conn.QueryRowContext(ctx, getQuery, id)

	var record auditproto.Record
	err := row.Scan(&record.ID, &record.RequestID, &record.Action, &record.CreatedAt,
		&record.Data, &record.OriginIP, &record.ResourceID, &record.ResourceType, &record.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return auditproto.Record{}, derror.ErrUserNotFound
		}

		return auditproto.Record{}, err
	}

	return record, err
}

func (as *AuditStorage) List(ctx context.Context, page pagination.Page, orderBy order.OrderBy) ([]auditproto.Record, error) {
	conn := session.GetDBConnection(ctx, as.Conn)

	query := listQuery + page.String(pagination.SQLStringer(defaultPageSize)) + orderBy.String(oderStringer)

	rows, err := conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []auditproto.Record
	for rows.Next() {
		var record auditproto.Record
		err := rows.Scan(&record.ID, &record.RequestID, &record.Action, &record.CreatedAt,
			&record.Data, &record.OriginIP, &record.ResourceID, &record.ResourceType, &record.UserID)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

func (as *AuditStorage) Count(ctx context.Context) (int, error) {
	row := as.Conn.QueryRowContext(ctx, countQuery)
	var count int
	err := row.Scan(&count)
	return count, err
}

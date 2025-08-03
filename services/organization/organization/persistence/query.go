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
	orgproto "github.com/kianooshaz/skeleton/services/organization/organization/proto"
)

type OrganizationStorage struct {
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

func (os *OrganizationStorage) Create(ctx context.Context, organization orgproto.Organization) error {
	conn := session.GetDBConnection(ctx, os.Conn)

	_, err := conn.ExecContext(ctx, createQuery, organization.ID, organization.CreatedAt)
	return err
}

func (os *OrganizationStorage) Get(ctx context.Context, id orgproto.OrganizationID) (orgproto.Organization, error) {
	conn := session.GetDBConnection(ctx, os.Conn)

	row := conn.QueryRowContext(ctx, getQuery, id)

	var organization orgproto.Organization
	err := row.Scan(&organization.ID, &organization.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return orgproto.Organization{}, derror.ErrOrganizationNotFound
		}

		return orgproto.Organization{}, err
	}

	return organization, err
}

func (os *OrganizationStorage) List(ctx context.Context, page pagination.Page,
	orderBy order.OrderBy) ([]orgproto.Organization, error) {
	conn := session.GetDBConnection(ctx, os.Conn)

	query := listQuery + page.String(pagination.SQLStringer(defaultPageSize)) + orderBy.String(oderStringer)

	rows, err := conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var organizations []orgproto.Organization
	for rows.Next() {
		var organization orgproto.Organization
		err := rows.Scan(&organization.ID, &organization.CreatedAt)
		if err != nil {
			return nil, err
		}
		organizations = append(organizations, organization)
	}

	return organizations, nil
}

func (os *OrganizationStorage) Count(ctx context.Context) (int, error) {
	row := os.Conn.QueryRowContext(ctx, countQuery)
	var count int
	err := row.Scan(&count)
	return count, err
}

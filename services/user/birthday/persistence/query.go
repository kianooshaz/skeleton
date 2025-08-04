package persistence

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"strings"

	"github.com/kianooshaz/skeleton/foundation/order"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	"github.com/kianooshaz/skeleton/foundation/session"
	birthdayproto "github.com/kianooshaz/skeleton/services/user/birthday/proto"
	userproto "github.com/kianooshaz/skeleton/services/user/user/proto"
)

// ListFilters represents the filters for listing birthdays.
type ListFilters struct {
	UserID     *userproto.UserID
	MinAge     *int
	MaxAge     *int
	BirthMonth *int
}

//go:embed queries/create.sql
var createQuery string

//go:embed queries/get.sql
var getQuery string

//go:embed queries/get_by_user_id.sql
var getByUserIDQuery string

//go:embed queries/update.sql
var updateQuery string

//go:embed queries/delete.sql
var deleteQuery string

//go:embed queries/list.sql
var listQuery string

//go:embed queries/count.sql
var countQuery string

//go:embed queries/exists_by_user_id.sql
var existsByUserIDQuery string

// BirthdayStorage handles database operations for birthdays.
type BirthdayStorage struct {
	Conn *sql.DB
}

// Create creates a new birthday record in the database.
func (s *BirthdayStorage) Create(ctx context.Context, birthday birthdayproto.Birthday) error {
	conn := session.GetDBConnection(ctx, s.Conn)

	_, err := conn.ExecContext(
		ctx,
		createQuery,
		birthday.ID,
		birthday.UserID,
		birthday.DateOfBirth,
		birthday.Age,
	)
	if err != nil {
		return fmt.Errorf("creating birthday record: %w", err)
	}

	return nil
}

// Get retrieves a birthday record by ID.
func (s *BirthdayStorage) Get(ctx context.Context, id birthdayproto.BirthdayID) (birthdayproto.Birthday, error) {
	conn := session.GetDBConnection(ctx, s.Conn)

	var birthday birthdayproto.Birthday
	err := conn.QueryRowContext(ctx, getQuery, id).Scan(
		&birthday.ID,
		&birthday.UserID,
		&birthday.DateOfBirth,
		&birthday.Age,
		&birthday.CreatedAt,
		&birthday.UpdatedAt,
	)
	if err != nil {
		return birthdayproto.Birthday{}, fmt.Errorf("getting birthday record: %w", err)
	}

	return birthday, nil
}

// GetByUserID retrieves a birthday record by user ID.
func (s *BirthdayStorage) GetByUserID(ctx context.Context, userID userproto.UserID) (birthdayproto.Birthday, error) {
	conn := session.GetDBConnection(ctx, s.Conn)

	var birthday birthdayproto.Birthday
	err := conn.QueryRowContext(ctx, getByUserIDQuery, userID).Scan(
		&birthday.ID,
		&birthday.UserID,
		&birthday.DateOfBirth,
		&birthday.Age,
		&birthday.CreatedAt,
		&birthday.UpdatedAt,
	)
	if err != nil {
		return birthdayproto.Birthday{}, fmt.Errorf("getting birthday record by user ID: %w", err)
	}

	return birthday, nil
}

// Update updates an existing birthday record.
func (s *BirthdayStorage) Update(ctx context.Context, birthday birthdayproto.Birthday) error {
	conn := session.GetDBConnection(ctx, s.Conn)

	_, err := conn.ExecContext(
		ctx,
		updateQuery,
		birthday.ID,
		birthday.DateOfBirth,
		birthday.Age,
	)
	if err != nil {
		return fmt.Errorf("updating birthday record: %w", err)
	}

	return nil
}

// Delete removes a birthday record from the database.
func (s *BirthdayStorage) Delete(ctx context.Context, id birthdayproto.BirthdayID) error {
	conn := session.GetDBConnection(ctx, s.Conn)

	_, err := conn.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("deleting birthday record: %w", err)
	}

	return nil
}

// List retrieves a paginated list of birthday records with optional filters.
func (s *BirthdayStorage) List(ctx context.Context, page pagination.Page, orderBy order.OrderBy, filters ListFilters) ([]birthdayproto.Birthday, error) {
	conn := session.GetDBConnection(ctx, s.Conn)

	// Build the query with filters.
	query := listQuery
	args := []interface{}{}
	whereConditions := []string{}
	argIndex := 1

	if filters.UserID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("user_id = $%d", argIndex))
		args = append(args, *filters.UserID)
		argIndex++
	}

	if filters.MinAge != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("age >= $%d", argIndex))
		args = append(args, *filters.MinAge)
		argIndex++
	}

	if filters.MaxAge != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("age <= $%d", argIndex))
		args = append(args, *filters.MaxAge)
		argIndex++
	}

	if filters.BirthMonth != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("EXTRACT(MONTH FROM date_of_birth) = $%d", argIndex))
		args = append(args, *filters.BirthMonth)
	}

	if len(whereConditions) > 0 {
		query += " WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Add ordering.
	query += " ORDER BY " + orderStringer(orderBy)

	// Add pagination using the pattern from the framework.
	offset := page.PageRows * page.PageNumber
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", page.PageRows, offset)

	rows, err := conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("listing birthday records: %w", err)
	}
	defer rows.Close()

	var birthdays []birthdayproto.Birthday
	for rows.Next() {
		var birthday birthdayproto.Birthday
		err := rows.Scan(
			&birthday.ID,
			&birthday.UserID,
			&birthday.DateOfBirth,
			&birthday.Age,
			&birthday.CreatedAt,
			&birthday.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning birthday record: %w", err)
		}
		birthdays = append(birthdays, birthday)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating birthday records: %w", err)
	}

	return birthdays, nil
}

// Count returns the total number of birthday records matching the filters.
func (s *BirthdayStorage) Count(ctx context.Context, filters ListFilters) (int, error) {
	conn := session.GetDBConnection(ctx, s.Conn)

	// Build the query with filters.
	query := countQuery
	args := []interface{}{}
	whereConditions := []string{}
	argIndex := 1

	if filters.UserID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("user_id = $%d", argIndex))
		args = append(args, *filters.UserID)
		argIndex++
	}

	if filters.MinAge != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("age >= $%d", argIndex))
		args = append(args, *filters.MinAge)
		argIndex++
	}

	if filters.MaxAge != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("age <= $%d", argIndex))
		args = append(args, *filters.MaxAge)
		argIndex++
	}

	if filters.BirthMonth != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("EXTRACT(MONTH FROM date_of_birth) = $%d", argIndex))
		args = append(args, *filters.BirthMonth)
	}

	if len(whereConditions) > 0 {
		query = strings.Replace(query, "FROM birthdays", "FROM birthdays WHERE "+strings.Join(whereConditions, " AND "), 1)
	}

	var count int
	err := conn.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting birthday records: %w", err)
	}

	return count, nil
}

// ExistsByUserID checks if a birthday record exists for the given user ID.
func (s *BirthdayStorage) ExistsByUserID(ctx context.Context, userID userproto.UserID) (bool, error) {
	conn := session.GetDBConnection(ctx, s.Conn)

	var exists bool
	err := conn.QueryRowContext(ctx, existsByUserIDQuery, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking birthday existence by user ID: %w", err)
	}

	return exists, nil
}

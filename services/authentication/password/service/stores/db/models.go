// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Password struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	PasswordHash string
	CreatedAt    pgtype.Timestamptz
	DeletedAt    pgtype.Timestamptz
}

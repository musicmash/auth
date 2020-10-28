// Code generated by sqlc. DO NOT EDIT.
// source: users.sql

package models

import (
	"context"
	"database/sql"
)

const ensureUserExists = `-- name: EnsureUserExists :exec
INSERT INTO users (name, photo)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
`

type EnsureUserExistsParams struct {
	Name  string         `json:"name"`
	Photo sql.NullString `json:"photo"`
}

func (q *Queries) EnsureUserExists(ctx context.Context, arg EnsureUserExistsParams) error {
	_, err := q.db.ExecContext(ctx, ensureUserExists, arg.Name, arg.Photo)
	return err
}

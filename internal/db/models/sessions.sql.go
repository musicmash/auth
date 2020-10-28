// Code generated by sqlc. DO NOT EDIT.
// source: sessions.sql

package models

import (
	"context"
)

const createSession = `-- name: CreateSession :exec
INSERT INTO sessions (user_name, value)
VALUES ($1, $2)
`

type CreateSessionParams struct {
	UserName string `json:"user_name"`
	Value    string `json:"value"`
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) error {
	_, err := q.db.ExecContext(ctx, createSession, arg.UserName, arg.Value)
	return err
}

const getSession = `-- name: GetSession :one
SELECT id, created_at, user_name, value FROM sessions
where value = $1
`

func (q *Queries) GetSession(ctx context.Context, value string) (Session, error) {
	row := q.db.QueryRowContext(ctx, getSession, value)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UserName,
		&i.Value,
	)
	return i, err
}

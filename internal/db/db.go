package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

type Mgr struct {
	conn *pgx.Conn
}

func New(dsn string) (*Mgr, error) {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable connect to database: %w", err)
	}

	return &Mgr{conn: conn}, nil
}

func (mgr *Mgr) Close() error {
	return mgr.conn.Close(context.Background())
}

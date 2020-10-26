package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Mgr struct {
	conn *pgxpool.Pool
}

func New(dsn string) (*Mgr, error) {
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database_url: %w", err)
	}

	conn, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable connect to database: %w", err)
	}

	return &Mgr{conn: conn}, nil
}

func (mgr *Mgr) Close() {
	mgr.conn.Close()
}

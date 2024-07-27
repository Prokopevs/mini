package db

import (
	"context"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type database struct {
	db *sqlx.DB
}

func NewDatabase(ctx context.Context, addr string) (*database, error) {
	d, err := sqlx.ConnectContext(ctx, "pgx", addr)
	if err != nil {
		return nil, err
	}

	return &database{
		db: d,
	}, nil
}

package db

import (
	"context"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type database struct {
	db *sqlx.DB
	logger *zap.SugaredLogger
}

func NewDatabase(ctx context.Context, logger *zap.SugaredLogger, addr string) (*database, error) {
	d, err := sqlx.ConnectContext(ctx, "pgx", addr)
	if err != nil {
		return nil, err
	}

	return &database{
		db: d,
		logger: logger,
	}, nil
}

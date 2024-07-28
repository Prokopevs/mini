package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type txKey byte

const txKeyCtx txKey = 1

func setTxCtx(ctx context.Context, txx *sqlx.Tx) context.Context {
  return context.WithValue(ctx, txKeyCtx, txx)
}

func getTxCtx(ctx context.Context) *sqlx.Tx {
  txx, _ := ctx.Value(txKeyCtx).(*sqlx.Tx)
  return txx
}

func (t *database) RunTx(ctx context.Context, f func(context.Context) error) error {
	txx, err := t.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}

	txCtx := setTxCtx(ctx, txx)
	err = f(txCtx)
	if err != nil {
		txErr := txx.Rollback()
		if txErr != nil {
			t.logger.Errorw("rollback tx", "err", txErr)
		}

		return err
	}

	return txx.Commit()
}

func (t *database) GetExtContext(ctx context.Context) sqlx.ExtContext {
	if txx := getTxCtx(ctx); txx != nil {
		return txx
	}

	return t.db
}

func (t *database) GetQueryerContext(ctx context.Context) sqlx.QueryerContext {
	if txx := getTxCtx(ctx); txx != nil {
		return txx
	}

	return t.db
}

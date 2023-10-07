package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"route256/loms/internal/app/storage/repositories/queries"
)

type pgTransactions struct{}

type txKey struct{}

// beginTx begins new transaction if there is no outside transaction in context, otherwise it
// begins pseudo nested transaction on outside transaction
func (t *pgTransactions) beginTx(ctx context.Context, dbPool *pgxpool.Pool) (pgx.Tx, error) {
	outsideTx := t.getTxFromContext(ctx)

	if outsideTx == nil {
		return dbPool.Begin(ctx)
	}

	return outsideTx.Begin(ctx)
}

func (t *pgTransactions) getTxFromContext(ctx context.Context) pgx.Tx {
	if tx, ok := ctx.Value(txKey{}).(*pgxpool.Tx); ok {
		return tx
	}

	return nil
}

func (t *pgTransactions) getQueriesFromContext(ctx context.Context, dbPool *pgxpool.Pool) queries.Querier {
	if tx, ok := ctx.Value(txKey{}).(*pgxpool.Tx); ok {
		return queries.New(tx)
	}

	return queries.New(dbPool)
}

type PgTransactor struct {
	dbPool *pgxpool.Pool
}

func NewPgTransactor(pool *pgxpool.Pool) *PgTransactor {
	return &PgTransactor{dbPool: pool}
}

func (pt *PgTransactor) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	tx, err := pt.dbPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Err(err).Msg("cannot rollback transaction")
		}
	}(tx, ctx)

	err = tFunc(pt.injectTx(ctx, tx))
	if err != nil {
		// if error, rollback
		errRollback := tx.Rollback(ctx)
		if errRollback != nil && !errors.Is(errRollback, pgx.ErrTxClosed) {
			log.Err(errRollback).Msg("cannot rollback transaction")
			return err
		}
	}
	// if no error, commit
	err = tx.Commit(ctx)
	if err != nil {
		log.Err(err).Msg("cannot commit transaction")
		return err
	}
	return nil
}

func (pt *PgTransactor) injectTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

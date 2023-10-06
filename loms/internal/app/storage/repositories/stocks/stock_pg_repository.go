package stocks

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"route256/loms/internal/app/models"
	"route256/loms/internal/app/services"
	"route256/loms/internal/app/storage/repositories/stocks/queries"
)

type stocksPgRepository struct {
	dbPool *pgxpool.Pool
}

type txKey struct{}

func NewStocksPgRepository(dbPool *pgxpool.Pool) services.StocksProvider {
	return &stocksPgRepository{
		dbPool: dbPool,
	}
}

func (s stocksPgRepository) Reserve(ctx context.Context, order models.Order) error {
	tx, err := s.beginTx(ctx)
	if err != nil {
		log.Err(err).Msg("cannot begin transaction for reserving stocks")
		return fmt.Errorf("cannot begin transaction for reserving stocks: %w", err)
	}

	q := queries.New(tx)

	for _, item := range order.Items {
		params := queries.ChangeReserveOfSkuByAmountParams{
			Count: int64(item.Count),
			Sku:   int64(item.Sku),
		}
		err := q.ChangeReserveOfSkuByAmount(ctx, params)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Err(err).Msg("cannot commit transaction for reserving stocks")
		return fmt.Errorf("cannot commit transaction for reserving stocks: %w", err)
	}

	return nil
}

func (s stocksPgRepository) ReserveRemove(ctx context.Context, order models.Order) error {
	tx, err := s.beginTx(ctx)
	if err != nil {
		log.Err(err).Msg("cannot begin transaction for removing stocks")
		return fmt.Errorf("cannot begin transaction for removing stocks: %w", err)
	}

	q := queries.New(tx)

	for _, item := range order.Items {
		params := queries.ChangeReserveOfSkuByAmountParams{
			Count: -int64(item.Count),
			Sku:   int64(item.Sku),
		}
		err := q.ChangeReserveOfSkuByAmount(ctx, params)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Err(err).Msg("cannot commit transaction for removing stocks")
		return fmt.Errorf("cannot commit transaction for removing stocks: %w", err)
	}

	return nil
}

func (s stocksPgRepository) ReserveCancel(ctx context.Context, order models.Order) error {
	return s.ReserveRemove(ctx, order)
}

func (s stocksPgRepository) GetBySku(ctx context.Context, sku uint32) (uint64, error) {
	q := s.getQueriesFromContext(ctx)
	countFromDb, err := q.GetBySku(ctx, int64(sku))
	if err != nil {
		return 0, err
	}

	return uint64(countFromDb), nil
}

// beginTx begins new transaction if there is no outside transaction in context, otherwise it
// begins pseudo nested transaction on outside transaction
func (s stocksPgRepository) beginTx(ctx context.Context) (pgx.Tx, error) {
	outsideTx := s.getTxFromContext(ctx)

	if outsideTx == nil {
		return s.dbPool.Begin(ctx)
	}

	return outsideTx.Begin(ctx)
}

func (s stocksPgRepository) getTxFromContext(ctx context.Context) pgx.Tx {
	if tx, ok := ctx.Value(txKey{}).(*pgxpool.Tx); ok {
		return tx
	}

	return nil
}

func (s stocksPgRepository) getQueriesFromContext(ctx context.Context) queries.Querier {
	if tx, ok := ctx.Value(txKey{}).(*pgxpool.Tx); ok {
		return queries.New(tx)
	}

	return queries.New(s.dbPool)
}

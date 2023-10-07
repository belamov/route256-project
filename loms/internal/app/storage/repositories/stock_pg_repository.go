package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"route256/loms/internal/app/storage"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"route256/loms/internal/app/models"
	"route256/loms/internal/app/storage/repositories/queries"
)

type StocksPgRepository struct {
	dbPool       *pgxpool.Pool
	transactions pgTransactions
}

func NewStocksPgRepository(dbPool *pgxpool.Pool) *StocksPgRepository {
	return &StocksPgRepository{
		dbPool:       dbPool,
		transactions: pgTransactions{},
	}
}

var CheckViolationErrorCode = "23514"

func (s StocksPgRepository) Reserve(ctx context.Context, order models.Order) error {
	tx, err := s.transactions.beginTx(ctx, s.dbPool)
	if err != nil {
		log.Err(err).Msg("cannot begin transaction for reserving stocks")
		return fmt.Errorf("cannot begin transaction for reserving stocks: %w", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Err(err).Msg("cannot rollback transaction")
		}
	}(tx, ctx)

	q := queries.New(tx)

	for _, item := range order.Items {
		params := queries.ChangeReserveOfSkuByAmountParams{
			Count: -int64(item.Count),
			Sku:   int64(item.Sku),
		}
		res, err := q.ChangeReserveOfSkuByAmount(ctx, params)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == CheckViolationErrorCode && strings.Contains(pgErr.Message, "count_nonnegative") {
					return storage.ErrInsufficientStocks
				}
				return err
			}
		}
		if res.RowsAffected() != 1 {
			return storage.ErrInsufficientStocks
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Err(err).Msg("cannot commit transaction for reserving stocks")
		return fmt.Errorf("cannot commit transaction for reserving stocks: %w", err)
	}

	return nil
}

func (s StocksPgRepository) ReserveRemove(ctx context.Context, order models.Order) error {
	tx, err := s.transactions.beginTx(ctx, s.dbPool)
	if err != nil {
		log.Err(err).Msg("cannot begin transaction for removing stocks")
		return fmt.Errorf("cannot begin transaction for removing stocks: %w", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Err(err).Msg("cannot rollback transaction")
		}
	}(tx, ctx)

	q := queries.New(tx)

	for _, item := range order.Items {
		params := queries.ChangeReserveOfSkuByAmountParams{
			Count: int64(item.Count),
			Sku:   int64(item.Sku),
		}
		res, err := q.ChangeReserveOfSkuByAmount(ctx, params)
		if err != nil {
			return err
		}
		if res.RowsAffected() != 1 {
			return errors.New("sku not found in stocks")
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Err(err).Msg("cannot commit transaction for removing stocks")
		return fmt.Errorf("cannot commit transaction for removing stocks: %w", err)
	}

	return nil
}

func (s StocksPgRepository) ReserveCancel(ctx context.Context, order models.Order) error {
	return s.ReserveRemove(ctx, order)
}

func (s StocksPgRepository) GetBySku(ctx context.Context, sku uint32) (uint64, error) {
	q := s.transactions.getQueriesFromContext(ctx, s.dbPool)
	countFromDb, err := q.GetBySku(ctx, int64(sku))
	if err != nil {
		return 0, err
	}

	return uint64(countFromDb), nil
}

package cart

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"route256/loms/internal/app/models"
	"route256/loms/internal/app/services"
	"route256/loms/internal/app/storage/repositories/stocks/queries"
	"sync"
)

type stocksPgRepository struct {
	q      *queries.Queries
	dbPool *pgxpool.Pool
}

func NewStocksRepository(ctx context.Context, wg *sync.WaitGroup, user string, password string, host string, db string) (services.StocksProvider, error) {
	databaseDSN := fmt.Sprintf("postgresql://%s:%s@%s/%s", user, password, host, db)
	dbPool, err := pgxpool.New(ctx, databaseDSN)
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		log.Info().Msg("Closing order repository connections...")
		dbPool.Close()
		log.Info().Msg("Order repository connections closed")
		wg.Done()
	}()

	return &stocksPgRepository{
		q:      queries.New(dbPool),
		dbPool: dbPool,
	}, nil
}

func (s stocksPgRepository) Reserve(ctx context.Context, order models.Order) error {
	tx, err := s.dbPool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Err(err).Msg("cannot rollback transaction for reserving stocks")
		}
	}(tx, ctx)

	q := s.q.WithTx(tx)
	for _, item := range order.Items {
		params := queries.ReserveSkuParams{
			Count: int64(item.Count),
			Sku:   int64(item.Sku),
		}
		err := q.ReserveSku(ctx, params)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s stocksPgRepository) ReserveRemove(ctx context.Context, order models.Order) error {
	tx, err := s.dbPool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Err(err).Msg("cannot rollback transaction for removing reserves")
		}
	}(tx, ctx)

	q := s.q.WithTx(tx)
	for _, item := range order.Items {
		params := queries.RemoveReserveSkuParams{
			Count: int64(item.Count),
			Sku:   int64(item.Sku),
		}
		err := q.RemoveReserveSku(ctx, params)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s stocksPgRepository) ReserveCancel(ctx context.Context, order models.Order) error {
	//TODO implement me
	panic("implement me")
}

func (s stocksPgRepository) GetBySku(ctx context.Context, sku uint32) (uint64, error) {
	countFromDb, err := s.q.GetBySku(ctx, int64(sku))
	if err != nil {
		return 0, err
	}

	return uint64(countFromDb), nil
}

package cart

import (
	"context"
	"fmt"
	"route256/loms/internal/app/models"
	"route256/loms/internal/app/services"
	"route256/loms/internal/app/storage/repositories/order/queries"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type orderPgRepository struct {
	q queries.Querier
}

func NewOrderRepository(ctx context.Context, wg *sync.WaitGroup, user string, password string, host string, db string) (services.OrdersProvider, error) {
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

	return &orderPgRepository{
		q: queries.New(dbPool),
	}, nil
}

func (o *orderPgRepository) Create(ctx context.Context, userId int64, statusNew models.OrderStatus, items []models.OrderItem) (models.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *orderPgRepository) SetStatus(ctx context.Context, order models.Order, status models.OrderStatus) (models.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *orderPgRepository) GetOrderByOrderId(ctx context.Context, orderId int64) (models.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *orderPgRepository) CancelUnpaidOrders(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (o *orderPgRepository) GetOrdersIdsByCreatedAtAndStatus(ctx context.Context, createdAt time.Time, orderStatus models.OrderStatus) ([]int64, error) {
	//TODO implement me
	panic("implement me")
}

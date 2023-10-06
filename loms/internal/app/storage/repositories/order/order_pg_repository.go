package order

import (
	"context"
	"route256/loms/internal/app/models"
	"route256/loms/internal/app/services"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type orderPgRepository struct {
	dbPool *pgxpool.Pool
}

func NewOrderRepository(dbPool *pgxpool.Pool) services.OrdersProvider {
	return &orderPgRepository{
		dbPool: dbPool,
	}
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

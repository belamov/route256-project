package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"route256/loms/internal/app/models"
	"route256/loms/internal/app/storage"
	"route256/loms/internal/app/storage/repositories/queries"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderPgRepository struct {
	dbPool       *pgxpool.Pool
	transactions pgTransactions
}

func NewOrderPgRepository(dbPool *pgxpool.Pool) *OrderPgRepository {
	return &OrderPgRepository{
		dbPool:       dbPool,
		transactions: pgTransactions{},
	}
}

func (o *OrderPgRepository) Create(ctx context.Context, userId int64, status models.OrderStatus, items []models.OrderItem) (models.Order, error) {
	tx, err := o.transactions.beginTx(ctx, o.dbPool)
	if err != nil {
		log.Err(err).Msg("cannot init transaction for creating order")
		return models.Order{}, err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Err(err).Msg("cannot rollback transaction")
		}
	}(tx, ctx)

	q := queries.New(tx)

	orderCreateParams := queries.CreateOrderParams{
		CreatedAt: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
		UserID: userId,
		Status: int16(status),
	}
	orderId, err := q.CreateOrder(ctx, orderCreateParams)
	if err != nil {
		log.Err(err).Msg("cannot create order")
		return models.Order{}, fmt.Errorf("cannot create order: %w", err)
	}

	for _, item := range items {
		createOrderItemParam := queries.CreateOrderItemsParams{
			OrderID: orderId,
			Name:    item.Name,
			Sku:     int64(item.Sku),
			Count:   int64(item.Count),
			Price:   int64(item.Price),
		}
		err = q.CreateOrderItems(ctx, createOrderItemParam)
		if err != nil {
			log.Err(err).Msg("cannot create order item")
			return models.Order{}, fmt.Errorf("cannot create order item: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Err(err).Msg("cannot commit order create transaction")
		return models.Order{}, fmt.Errorf("cannot commit order create transaction: %w", err)
	}

	order := models.Order{
		CreatedAt: orderCreateParams.CreatedAt.Time,
		Items:     items,
		Id:        orderId,
		UserId:    userId,
		Status:    status,
	}

	return order, nil
}

func (o *OrderPgRepository) SetStatus(ctx context.Context, order models.Order, status models.OrderStatus) (models.Order, error) {
	q := o.transactions.getQueriesFromContext(ctx, o.dbPool)

	params := queries.UpdateOrderStatusParams{
		Status: int16(status),
		ID:     order.Id,
	}
	err := q.UpdateOrderStatus(ctx, params)
	if err != nil {
		log.Err(err).Msg("cannot set status of order")
		return models.Order{}, fmt.Errorf("cannot set status of order: %w", err)
	}

	order.Status = status

	return order, nil
}

func (o *OrderPgRepository) GetOrderByOrderId(ctx context.Context, orderId int64) (models.Order, error) {
	q := o.transactions.getQueriesFromContext(ctx, o.dbPool)

	orderWithItems, err := q.GetOrderById(ctx, orderId)
	if err != nil {
		log.Err(err).Msg("cannot get order")
		return models.Order{}, fmt.Errorf("cannot get order: %w", err)
	}

	if len(orderWithItems) == 0 {
		return models.Order{}, storage.ErrOrderNotFound
	}

	orderFromDb := orderWithItems[0].Order
	order := models.Order{
		CreatedAt: orderFromDb.CreatedAt.Time,
		Items:     make([]models.OrderItem, 0, len(orderWithItems)),
		Id:        orderFromDb.ID,
		UserId:    orderFromDb.UserID,
		Status:    models.OrderStatus(orderFromDb.Status),
	}
	for _, orderItemFromDb := range orderWithItems {
		orderItem := models.OrderItem{
			Name:  orderItemFromDb.OrderItem.Name,
			User:  order.UserId,
			Sku:   uint32(orderItemFromDb.OrderItem.Sku),
			Price: uint32(orderItemFromDb.OrderItem.Price),
			Count: uint64(orderItemFromDb.OrderItem.Count),
		}
		order.Items = append(order.Items, orderItem)
	}

	return order, nil
}

func (o *OrderPgRepository) GetExpiredOrdersWithStatus(ctx context.Context, createdAt time.Time, orderStatus models.OrderStatus) ([]int64, error) {
	q := o.transactions.getQueriesFromContext(ctx, o.dbPool)

	params := queries.GetExpiredOrdersWithStatusParams{
		CreatedAt: pgtype.Timestamp{Time: createdAt, Valid: true},
		Status:    int16(orderStatus),
	}
	ordersIds, err := q.GetExpiredOrdersWithStatus(ctx, params)
	if err != nil {
		log.Err(err).Msg("cannot get expired orders ids")
		return nil, fmt.Errorf("cannot get expired orders ids: %w", err)
	}

	return ordersIds, nil
}

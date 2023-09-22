package services

import (
	"context"
	"errors"
	"fmt"

	"route256/loms/internal/app/storage"

	"github.com/rs/zerolog/log"
	"route256/loms/internal/app/models"
)

// TODO:
// - отмена заказа после 10 минут неуплаты
// - валидация смены статусов заказа
type Loms interface {
	OrderCreate(ctx context.Context, userId int64, items []models.OrderItem) (models.Order, error)
	OrderInfo(ctx context.Context, orderId int64) (models.Order, error)
	OrderPay(ctx context.Context, orderId int64) error
	OrderCancel(ctx context.Context, orderId int64) error
	StockInfo(ctx context.Context, sku uint32) (uint64, error)
}

var (
	ErrInsufficientStocks = errors.New("insufficient stocks")
	ErrOrderNotFound      = errors.New("order not found")
)

type OrdersProvider interface {
	Create(ctx context.Context, userId int64, statusNew models.OrderStatus, items []models.OrderItem) (models.Order, error)
	SetStatus(ctx context.Context, order models.Order, status models.OrderStatus) (models.Order, error)
	GetOrderByOrderId(ctx context.Context, orderId int64) (models.Order, error)
}

type StocksProvider interface {
	Reserve(ctx context.Context, order models.Order) error
	ReserveRemove(ctx context.Context, order models.Order) error
	ReserveCancel(ctx context.Context, order models.Order) error
	GetBySku(ctx context.Context, sku uint32) (uint64, error)
}

type lomsService struct {
	ordersProvider OrdersProvider
	stocksProvider StocksProvider
}

func NewLomsService(
	ordersProvider OrdersProvider,
	stocksProvider StocksProvider,
) Loms {
	return &lomsService{
		ordersProvider: ordersProvider,
		stocksProvider: stocksProvider,
	}
}

func (l *lomsService) OrderCreate(ctx context.Context, userId int64, items []models.OrderItem) (models.Order, error) {
	order, err := l.ordersProvider.Create(ctx, userId, models.OrderStatusNew, items)
	if err != nil {
		log.Err(err).
			Msg("failed creating new order!")
		return models.Order{}, fmt.Errorf("failed creating new order!: %w", err)
	}

	err = l.stocksProvider.Reserve(ctx, order)
	if err != nil {
		failedOrder, errSetStatus := l.ordersProvider.SetStatus(ctx, order, models.OrderStatusFailed)
		if errSetStatus != nil {
			log.Err(errSetStatus).
				Int64("orderId", failedOrder.Id).
				Msg("failed transition order to failed status!")
			return models.Order{}, fmt.Errorf("failed transition order to failed status!: %w", errSetStatus)
		}

		return models.Order{}, ErrInsufficientStocks
	}

	awaitingOrder, err := l.ordersProvider.SetStatus(ctx, order, models.OrderStatusAwaitingPayment)
	if err != nil {
		log.Err(err).
			Int64("orderId", awaitingOrder.Id).
			Msg("failed transition order to awaiting status!")
		return models.Order{}, fmt.Errorf("failed transition order to failed status!: %w", err)
	}

	return awaitingOrder, nil
}

func (l *lomsService) OrderInfo(ctx context.Context, orderId int64) (models.Order, error) {
	order, err := l.ordersProvider.GetOrderByOrderId(ctx, orderId)
	if errors.Is(err, storage.ErrOrderNotFound) {
		return models.Order{}, ErrOrderNotFound
	}
	if err != nil {
		log.Err(err).
			Int64("orderId", orderId).
			Msg("failed getting order!")
		return models.Order{}, fmt.Errorf("failed getting order!: %w", err)
	}

	return order, nil
}

func (l *lomsService) OrderPay(ctx context.Context, orderId int64) error {
	order, err := l.ordersProvider.GetOrderByOrderId(ctx, orderId)
	if errors.Is(err, storage.ErrOrderNotFound) {
		return ErrOrderNotFound
	}
	if err != nil {
		log.Err(err).
			Int64("orderId", orderId).
			Msg("failed getting order!")
		return fmt.Errorf("failed getting order!: %w", err)
	}

	err = l.stocksProvider.ReserveRemove(ctx, order)
	if err != nil {
		log.Err(err).
			Any("orderId", orderId).
			Msg("failed removing reserves!")
		return fmt.Errorf("failed removing reserve!: %w", err)
	}

	_, err = l.ordersProvider.SetStatus(ctx, order, models.OrderStatusPayed)
	if err != nil {
		log.Err(err).
			Any("orderId", orderId).
			Msg("failed setting order status to payed!")
		return fmt.Errorf("failed setting order status to payed!: %w", err)
	}

	return nil
}

func (l *lomsService) OrderCancel(ctx context.Context, orderId int64) error {
	order, err := l.ordersProvider.GetOrderByOrderId(ctx, orderId)
	if errors.Is(err, storage.ErrOrderNotFound) {
		return ErrOrderNotFound
	}
	if err != nil {
		log.Err(err).
			Int64("orderId", orderId).
			Msg("failed getting order!")
		return fmt.Errorf("failed getting order!: %w", err)
	}

	err = l.stocksProvider.ReserveCancel(ctx, order)
	if err != nil {
		log.Err(err).
			Any("orderId", orderId).
			Msg("failed canceling reserves!")
		return fmt.Errorf("failed canceling reserve!: %w", err)
	}

	_, err = l.ordersProvider.SetStatus(ctx, order, models.OrderStatusCancelled)
	if err != nil {
		log.Err(err).
			Any("orderId", orderId).
			Msg("failed setting order status to canceled!")
		return fmt.Errorf("failed setting order status to canceled!: %w", err)
	}

	return nil
}

func (l *lomsService) StockInfo(ctx context.Context, sku uint32) (uint64, error) {
	count, err := l.stocksProvider.GetBySku(ctx, sku)
	if err != nil {
		return 0, err
	}

	return count, nil
}

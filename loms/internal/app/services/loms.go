package services

import (
	"context"
	"errors"
	"fmt"

	"route256/loms/internal/app/storage"

	"github.com/rs/zerolog/log"
	"route256/loms/internal/app/models"
)

type Loms interface {
	OrderCreate(ctx context.Context, userId int64, items []models.OrderItem) (models.Order, error)
	OrderInfo(ctx context.Context, orderId int64) (models.Order, error)
	OrderPay(ctx context.Context, id int64) error
}

var (
	ErrInsufficientStocks = errors.New("insufficient stocks")
	ErrOrderNotFound      = errors.New("order not found")
)

type ProductService interface{}

type OrdersStorage interface {
	Create(ctx context.Context, userId int64, statusNew models.OrderStatus, items []models.OrderItem) (models.Order, error)
	SetStatus(ctx context.Context, order models.Order, status models.OrderStatus) (models.Order, error)
	GetOrderByOrderId(ctx context.Context, orderId int64) (models.Order, error)
}

type StocksStorage interface {
	Reserve(ctx context.Context, order models.Order) error
	ReserveRemove(ctx context.Context, order models.Order) error
}

type lomsService struct {
	productService ProductService
	ordersStorage  OrdersStorage
	stocksStorage  StocksStorage
}

func NewLomsService(
	productService ProductService,
	ordersStorage OrdersStorage,
	stocksStorage StocksStorage,
) Loms {
	return &lomsService{
		productService: productService,
		ordersStorage:  ordersStorage,
		stocksStorage:  stocksStorage,
	}
}

func (l *lomsService) OrderCreate(ctx context.Context, userId int64, items []models.OrderItem) (models.Order, error) {
	order, err := l.ordersStorage.Create(ctx, userId, models.OrderStatusNew, items)
	if err != nil {
		log.Err(err).
			Msg("failed creating new order!")
		return models.Order{}, fmt.Errorf("failed creating new order!: %w", err)
	}

	err = l.stocksStorage.Reserve(ctx, order)
	if err != nil {
		failedOrder, errSetStatus := l.ordersStorage.SetStatus(ctx, order, models.OrderStatusFailed)
		if errSetStatus != nil {
			log.Err(errSetStatus).
				Int64("orderId", failedOrder.Id).
				Msg("failed transition order to failed status!")
			return models.Order{}, fmt.Errorf("failed transition order to failed status!: %w", errSetStatus)
		}

		return models.Order{}, ErrInsufficientStocks
	}

	awaitingOrder, err := l.ordersStorage.SetStatus(ctx, order, models.OrderStatusAwaitingPayment)
	if err != nil {
		log.Err(err).
			Int64("orderId", awaitingOrder.Id).
			Msg("failed transition order to awaiting status!")
		return models.Order{}, fmt.Errorf("failed transition order to failed status!: %w", err)
	}

	return awaitingOrder, nil
}

func (l *lomsService) OrderInfo(ctx context.Context, orderId int64) (models.Order, error) {
	order, err := l.ordersStorage.GetOrderByOrderId(ctx, orderId)
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
	order, err := l.ordersStorage.GetOrderByOrderId(ctx, orderId)
	if errors.Is(err, storage.ErrOrderNotFound) {
		return ErrOrderNotFound
	}
	if err != nil {
		log.Err(err).
			Int64("orderId", orderId).
			Msg("failed getting order!")
		return fmt.Errorf("failed getting order!: %w", err)
	}

	err = l.stocksStorage.ReserveRemove(ctx, order)
	if err != nil {
		log.Err(err).
			Any("orderId", orderId).
			Msg("failed removing reserves!")
		return fmt.Errorf("failed removing reserve!: %w", err)
	}

	_, err = l.ordersStorage.SetStatus(ctx, order, models.OrderStatusPayed)
	if err != nil {
		log.Err(err).
			Any("orderId", orderId).
			Msg("failed setting order status to payed!")
		return fmt.Errorf("failed setting order status to payed!: %w", err)
	}

	return nil
}

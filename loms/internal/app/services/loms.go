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
	CreateOrder(ctx context.Context, userId int64, items []models.OrderItem) (models.Order, error)
	GetOrderById(ctx context.Context, orderId int64) (models.Order, error)
}

var (
	ErrInsufficientStocks = errors.New("insufficient stocks")
	ErrOrderNotFound      = errors.New("order not found")
)

type ProductService interface{}

type OrdersStorage interface {
	Create(ctx context.Context, userId int64, statusNew models.OrderStatus, items []models.OrderItem) (models.Order, error)
	SetStatus(ctx context.Context, order models.Order, status models.OrderStatus) (models.Order, error)
	GetById(ctx context.Context, order int64) (models.Order, error)
}

type StocksStorage interface {
	Reserve(ctx context.Context, items []models.OrderItem) error
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

func (l *lomsService) CreateOrder(ctx context.Context, userId int64, items []models.OrderItem) (models.Order, error) {
	order, err := l.ordersStorage.Create(ctx, userId, models.OrderStatusNew, items)
	if err != nil {
		log.Err(err).
			Msg("failed creating new order!")
		return models.Order{}, fmt.Errorf("failed creating new order!: %w", err)
	}

	err = l.stocksStorage.Reserve(ctx, items)
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

func (l *lomsService) GetOrderById(ctx context.Context, orderId int64) (models.Order, error) {
	order, err := l.ordersStorage.GetById(ctx, orderId)
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

package services

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"route256/loms/internal/app/models"

	"route256/loms/internal/app/storage"

	"github.com/rs/zerolog/log"
)

type Loms interface {
	OrderCreate(ctx context.Context, userId int64, items []models.OrderItem) (models.Order, error)
	OrderInfo(ctx context.Context, orderId int64) (models.Order, error)
	OrderPay(ctx context.Context, orderId int64) error
	OrderCancel(ctx context.Context, orderId int64) error
	StockInfo(ctx context.Context, sku uint32) (uint64, error)
	RunCancelUnpaidOrders(ctx context.Context, wg *sync.WaitGroup, period time.Duration)
}

var (
	ErrInsufficientStocks = errors.New("insufficient stocks")
	ErrOrderNotFound      = errors.New("order not found")
	ErrOrderCancelled     = errors.New("order canceled")
)

type OrdersProvider interface {
	Create(ctx context.Context, userId int64, statusNew models.OrderStatus, items []models.OrderItem) (models.Order, error)
	SetStatus(ctx context.Context, order models.Order, status models.OrderStatus) (models.Order, error)
	GetOrderByOrderId(ctx context.Context, orderId int64) (models.Order, error)
	GetExpiredOrdersWithStatus(ctx context.Context, createdAt time.Time, orderStatus models.OrderStatus) ([]int64, error)
}

type StocksProvider interface {
	Reserve(ctx context.Context, order models.Order) error
	ReserveRemove(ctx context.Context, order models.Order) error
	ReserveCancel(ctx context.Context, order models.Order) error
	GetBySku(ctx context.Context, sku uint32) (uint64, error)
}

type Transactor interface {
	WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error
}

type OrderEventsProducer interface {
	OrderStatusChangedEventEmit(ctx context.Context, order models.Order) error
}

type lomsService struct {
	ordersProvider         OrdersProvider
	stocksProvider         StocksProvider
	allowedOrderUnpaidTime time.Duration
	transactor             Transactor
	orderEventsProducer    OrderEventsProducer
}

const DefaultAllowedOrderUnpaidTime = time.Minute * 10

func NewLomsService(
	ordersProvider OrdersProvider,
	stocksProvider StocksProvider,
	allowedOrderUnpaidTime time.Duration,
	transactor Transactor,
	orderEventsProducer OrderEventsProducer,
) Loms {
	if allowedOrderUnpaidTime == 0 {
		allowedOrderUnpaidTime = DefaultAllowedOrderUnpaidTime
	}
	return &lomsService{
		ordersProvider:         ordersProvider,
		stocksProvider:         stocksProvider,
		allowedOrderUnpaidTime: allowedOrderUnpaidTime,
		transactor:             transactor,
		orderEventsProducer:    orderEventsProducer,
	}
}

func (l *lomsService) OrderCreate(ctx context.Context, userId int64, items []models.OrderItem) (models.Order, error) {
	order, err := l.ordersProvider.Create(ctx, userId, models.OrderStatusNew, items)
	if err != nil {
		log.Err(err).
			Msg("failed creating new order!")
		return models.Order{}, fmt.Errorf("failed creating new order!: %w", err)
	}

	err = l.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		err = l.stocksProvider.Reserve(ctx, order)
		if errors.Is(err, storage.ErrInsufficientStocks) {
			failedOrder, errSetStatus := l.SetOrderStatus(ctx, order, models.OrderStatusFailed)
			if errSetStatus != nil {
				log.Err(errSetStatus).
					Int64("orderId", failedOrder.Id).
					Msg("failed transition order to failed status!")
				return fmt.Errorf("failed transition order to failed status!: %w", errSetStatus)
			}

			return ErrInsufficientStocks
		}
		if err != nil {
			return fmt.Errorf("cannot reserve stocks: %w", err)
		}

		order, err = l.SetOrderStatus(ctx, order, models.OrderStatusAwaitingPayment)
		if err != nil {
			log.Err(err).
				Int64("orderId", order.Id).
				Msg("failed transition order to awaiting status!")
			return fmt.Errorf("failed transition order to failed status!: %w", err)
		}

		return nil
	})

	return order, err
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

	// у нас есть авто отмена неоплаченных заказов спустя
	// время AllowedOrderUnpaidTime после формирования заказа (см. RunCancelUnpaidOrders)
	// чтобы избежать ситуаций, когда время отмены заказа совпадет с его оплатой, дополнительно проверим, что
	// заказ не будет отменен
	if order.ShouldBeCancelled(l.allowedOrderUnpaidTime) {
		return ErrOrderCancelled
	}

	return l.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		err = l.stocksProvider.ReserveRemove(ctx, order)
		if err != nil {
			log.Err(err).
				Any("orderId", orderId).
				Msg("failed removing reserves!")
			return fmt.Errorf("failed removing reserve!: %w", err)
		}

		_, err = l.SetOrderStatus(ctx, order, models.OrderStatusPayed)
		if err != nil {
			log.Err(err).
				Any("orderId", orderId).
				Msg("failed setting order status to payed!")
			return fmt.Errorf("failed setting order status to payed!: %w", err)
		}

		return nil
	})
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

	return l.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		err = l.stocksProvider.ReserveCancel(ctx, order)
		if err != nil {
			log.Err(err).
				Any("orderId", orderId).
				Msg("failed canceling reserves!")
			return fmt.Errorf("failed canceling reserve!: %w", err)
		}

		_, err = l.SetOrderStatus(ctx, order, models.OrderStatusCancelled)
		if err != nil {
			log.Err(err).
				Any("orderId", orderId).
				Msg("failed setting order status to canceled!")
			return fmt.Errorf("failed setting order status to canceled!: %w", err)
		}

		return nil
	})
}

func (l *lomsService) StockInfo(ctx context.Context, sku uint32) (uint64, error) {
	count, err := l.stocksProvider.GetBySku(ctx, sku)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (l *lomsService) RunCancelUnpaidOrders(ctx context.Context, wg *sync.WaitGroup, period time.Duration) {
	ticker := time.NewTicker(period)
	log.Info().Msg("Starting canceling unpaid orders")
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			log.Info().Msg("Stopped canceling unpaid orders")
			wg.Done()
			return

		case <-ticker.C:
			log.Info().Msg("Cancelling unpaid orders")
			// дополнительно отнимем минуту, чтобы предотвратить конфликты. при оплате заказа проверяется
			// время без этой минуты, а значит будет некий буфер, для того, чтобы точно не оплачивать отмененные заказы
			timeUnpaidOrdersShouldBeCancelled := time.Now().Add(-l.allowedOrderUnpaidTime - time.Minute)
			ordersIds, err := l.ordersProvider.GetExpiredOrdersWithStatus(
				ctx,
				timeUnpaidOrdersShouldBeCancelled,
				models.OrderStatusAwaitingPayment,
			)
			if err != nil {
				log.Err(err).Msg("failed to fetch orders to cancel")
			}

			log.Info().Ints64("ordersIds", ordersIds).Msg("Cancelling unpaid orders")
			for _, orderId := range ordersIds {
				err = l.OrderCancel(context.Background(), orderId)
				if err != nil {
					log.Err(err).Msg("failed to cancel unpaid order")
				}
			}
		}
	}
}

func (l *lomsService) SetOrderStatus(ctx context.Context, order models.Order, newStatus models.OrderStatus) (models.Order, error) {
	updatedOrder, err := l.ordersProvider.SetStatus(ctx, order, newStatus)
	if err != nil {
		return updatedOrder, fmt.Errorf("failed to set order status: %w", err)
	}

	err = l.orderEventsProducer.OrderStatusChangedEventEmit(ctx, updatedOrder)
	if err != nil {
		return updatedOrder, fmt.Errorf("failed to emit order status changed event: %w", err)
	}

	return updatedOrder, nil
}

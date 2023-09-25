package models

import "time"

type OrderStatus uint16

const (
	OrderStatusNew OrderStatus = iota
	OrderStatusAwaitingPayment
	OrderStatusFailed
	OrderStatusPayed
	OrderStatusCancelled
)

func (s OrderStatus) String() string {
	switch s {
	case OrderStatusNew:
		return "new"
	case OrderStatusAwaitingPayment:
		return "awaiting payment"
	case OrderStatusFailed:
		return "failed"
	case OrderStatusPayed:
		return "payed"
	case OrderStatusCancelled:
		return "canceled"
	}
	return "unknown"
}

type Order struct {
	CreatedAt time.Time
	Items     []OrderItem
	Id        int64
	UserId    int64
	Status    OrderStatus
}

func (o *Order) ShouldBeCancelled(allowedOrderUnpaidTime time.Duration) bool {
	if o.Status != OrderStatusAwaitingPayment {
		return false
	}

	durationOrderUnpaid := time.Since(o.CreatedAt)
	return durationOrderUnpaid >= allowedOrderUnpaidTime
}

type OrderItem struct {
	Sku   int32
	Count uint16
}

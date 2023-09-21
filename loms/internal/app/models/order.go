package models

type OrderStatus uint16

const (
	OrderStatusUndefined OrderStatus = iota
	OrderStatusNew
	OrderStatusFailed
	OrderStatusAwaitingPayment
)

func (s OrderStatus) String() string {
	switch s {
	case OrderStatusUndefined:
		return "undefined"
	case OrderStatusNew:
		return "new"
	case OrderStatusFailed:
		return "failed"
	case OrderStatusAwaitingPayment:
		return "awaiting_payment"
	}
	return "unknown"
}

type Order struct {
	Items  []OrderItem
	Id     int64
	UserId int64
	Status OrderStatus
}

type OrderItem struct {
	Sku   int32
	Count uint16
}

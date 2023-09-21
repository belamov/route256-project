package models

type OrderStatus uint16

const (
	OrderStatusUndefined OrderStatus = iota
	OrderStatusNew
	OrderStatusFailed
	OrderStatusAwaitingPayment
)

type Order struct {
	Items  []OrderItem
	Id     int
	Status OrderStatus
}

type OrderItem struct {
	Sku   int32
	Count uint16
}

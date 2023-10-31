package models

type OutboxMessage struct {
	Id          int64
	Key         string
	Destination string
	Data        []byte
}

type OutboxFailedMessage struct {
	Message OutboxMessage
	Error   error
}

type OrderStatusInfo struct {
	OrderId   int64       `json:"order_id"`
	NewStatus OrderStatus `json:"new_status"`
}

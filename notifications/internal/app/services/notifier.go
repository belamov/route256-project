package services

import "github.com/rs/zerolog/log"

type Notifier struct{}

func NewNotifier() *Notifier {
	return &Notifier{}
}

func (n *Notifier) NotifyAboutOrderStatusChange(orderId int64, newOrderStatus uint16) error {
	log.Info().
		Int64("order id", orderId).
		Uint16("new order status", newOrderStatus).
		Msg("order status changed")
	return nil
}

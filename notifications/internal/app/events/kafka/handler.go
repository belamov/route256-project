package kafka

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/IBM/sarama"
)

type Notifier interface {
	NotifyAboutOrderStatusChange(orderId int64, newOrderStatus uint16) error
}

type orderStatusInfo struct {
	OrderId   int64  `json:"order_id"`
	NewStatus uint16 `json:"new_status"`
}

func (info *orderStatusInfo) validate() error {
	if info.OrderId == 0 {
		return errors.New("invalid message: order_id field is required")
	}
	if info.NewStatus == 0 {
		return errors.New("invalid message: new status is required")
	}
	if info.NewStatus > 4 {
		return errors.New("invalid message: status is more than 4")
	}
	return nil
}

var _ sarama.ConsumerGroupHandler = (*orderStatusChangedHandler)(nil)

type orderStatusChangedHandler struct {
	notifier Notifier
}

func newOrderStatusChangedHandler(notifier Notifier) *orderStatusChangedHandler {
	return &orderStatusChangedHandler{notifier: notifier}
}

func (h *orderStatusChangedHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *orderStatusChangedHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *orderStatusChangedHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():

			var data orderStatusInfo
			err := json.Unmarshal(message.Value, &data)
			if err != nil {
				session.MarkMessage(message, "")
				return fmt.Errorf("failed unmarshallng order status info: %w", err)
			}

			err = data.validate()
			if err != nil {
				session.MarkMessage(message, "")
				return fmt.Errorf("received invalid message: %w", err)
			}

			err = h.notifier.NotifyAboutOrderStatusChange(data.OrderId, data.NewStatus)
			if err != nil {
				session.MarkMessage(message, "")
				return fmt.Errorf("error from notify service: %w", err)
			}
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

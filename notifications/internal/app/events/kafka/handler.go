package kafka

import (
	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)

var _ sarama.ConsumerGroupHandler = (*orderStatusChangedHandler)(nil)

type orderStatusChangedHandler struct{}

func newOrderStatusChangedHandler() *orderStatusChangedHandler {
	return &orderStatusChangedHandler{}
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
			h.handleOrderStatusChange(message)

			// коммит сообщения "руками"
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

func (h *orderStatusChangedHandler) handleOrderStatusChange(message *sarama.ConsumerMessage) {
	log.Info().Msg("new event from kafka: " + string(message.Value))
}

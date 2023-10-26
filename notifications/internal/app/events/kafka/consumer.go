package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
)

type Consumer struct {
	consumerGroup *consumerGroup
}

func NewConsumer(ctx context.Context, wg *sync.WaitGroup, brokers []string, topicNames []string, consumerGroupId string) (*Consumer, error) {
	// Наш обработчик реализующий интерфейс sarama.orderStatusChangedHandler
	handler := newOrderStatusChangedHandler()

	// Создаем коньюмер группу
	cg, err := newConsumerGroup(
		brokers,
		consumerGroupId,
		topicNames,
		handler,
	)
	if err != nil {
		return nil, fmt.Errorf("error initializing kafka consumer group: %w", err)
	}

	log.Info().Msg("kafka consumer group configured. ready to consume messages")

	go func() {
		defer wg.Done()
		<-ctx.Done()
		log.Info().Msg("closing kafka consumer group...")
		if err = cg.Close(); err != nil {
			log.Error().Err(err).Msg("failed closing kafka consumer group")
			return
		}
		log.Info().Msg("closed kafka consumer group")
	}()

	return &Consumer{
		consumerGroup: cg,
	}, nil
}

func (k Consumer) StartConsuming(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Info().Msg("starting consuming events")
	// запускаем вычитку сообщений
	k.consumerGroup.Run(ctx)
}

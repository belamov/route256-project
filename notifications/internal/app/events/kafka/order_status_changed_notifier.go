package kafka

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/IBM/sarama"

	"github.com/rs/zerolog/log"
)

type OrderStatusChangedNotifier struct {
	consumerGroup sarama.ConsumerGroup
	handler       sarama.ConsumerGroupHandler
	topics        []string
}

type Config struct {
	Brokers         []string
	TopicNames      []string
	ConsumerGroupId string
}

func NewOrderStatusChangedNotifier(
	ctx context.Context,
	wg *sync.WaitGroup,
	config Config,
	notifier Notifier,
) (*OrderStatusChangedNotifier, error) {
	// Наш обработчик реализующий интерфейс sarama.orderStatusChangedHandler
	handler := newOrderStatusChangedHandler(notifier)

	// Создаем коньюмер группу
	cg, err := newConsumerGroup(
		config.Brokers,
		config.ConsumerGroupId,
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

	return &OrderStatusChangedNotifier{
		consumerGroup: cg,
		handler:       handler,
		topics:        config.TopicNames,
	}, nil
}

func (c *OrderStatusChangedNotifier) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Info().Msg("starting consuming events")

	// Обязательно в случае config.Consumer.Return.Errors = true
	go func() {
		for err := range c.consumerGroup.Errors() {
			log.Error().Err(err).Msg("Error from consumer")
		}
	}()

	for {
		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		if err := c.consumerGroup.Consume(ctx, c.topics, c.handler); err != nil {
			log.Error().Err(err).Msg("Error consuming message")
		}

		// check if context was cancelled, signaling that the consumer should stop
		if ctx.Err() != nil {
			return
		}
	}
}

func newConsumerGroup(brokers []string, groupID string) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion
	/*
		sarama.OffsetNewest - получаем только новые сообщений, те, которые уже были игнорируются
		sarama.OffsetOldest - читаем все с самого начала
	*/
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	// Используется, если ваш offset "уехал" далеко и нужно пропустить невалидные сдвиги
	config.Consumer.Group.ResetInvalidOffsets = true
	// Сердцебиение консьюмера
	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	// Таймаут сессии
	config.Consumer.Group.Session.Timeout = 60 * time.Second
	// Таймаут ребалансировки
	config.Consumer.Group.Rebalance.Timeout = 60 * time.Second
	//
	config.Consumer.Return.Errors = true

	const BalanceStrategy = "roundrobin"
	switch BalanceStrategy {
	case "sticky":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	case "roundrobin":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	case "range":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	default:
		return nil, fmt.Errorf("unrecognized consumer group partition assignor: %s", BalanceStrategy)
	}

	/*
	  Setup a new Sarama consumer group
	*/
	cg, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return cg, nil
}

package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"route256/loms/internal/app/models"
)

var OrderStatusChangedTopicName = "order-status-changed"

type Producer struct {
	producer sarama.AsyncProducer
}

func NewKafkaEventProducer(ctx context.Context, wg *sync.WaitGroup, brokers []string, opts ...Option) (*Producer, error) {
	config, err := prepareProducerSaramaConfig(opts...)
	if err != nil {
		return nil, err
	}

	asyncProducer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, errors.Wrap(err, "error with async kafka-producer")
	}

	go func() {
		defer wg.Done()
		<-ctx.Done()
		err := asyncProducer.Close()
		if err != nil {
			log.Err(err).Msg("error closing kafka producer")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Error и Retry топики можно использовать при получении ошибки
		for err := range asyncProducer.Errors() {
			log.Err(err).Msg("kafka async producer error")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range asyncProducer.Successes() {
			log.Info().Msg(fmt.Sprintf("Async success with key %s", msg.Key))
		}
	}()

	return &Producer{
		producer: asyncProducer,
	}, nil
}

type orderStatusInfo struct {
	OrderId   int64              `json:"order_id"`
	NewStatus models.OrderStatus `json:"new_status"`
}

func (p *Producer) OrderStatusChangedEventEmit(ctx context.Context, order models.Order) {
	orderInfo := orderStatusInfo{
		OrderId:   order.Id,
		NewStatus: order.Status,
	}

	bytes, err := json.Marshal(orderInfo)
	if err != nil {
		log.Err(err).Msg("failed to marshal order info")
	}

	msg, err := p.BuildMessage(OrderStatusChangedTopicName, fmt.Sprint(order.Id), bytes, "x-header-example", "example-header-value")
	if err != nil {
		log.Err(err).Msg("failed to build kafka message")
	}

	p.producer.Input() <- msg
}

func (p *Producer) BuildMessage(topic string, key string, message []byte, headersKV ...string) (*sarama.ProducerMessage, error) {
	if len(headersKV)%2 != 0 {
		return nil, errors.New("wrong number of headersKV")
	}

	headers := make([]sarama.RecordHeader, 0, len(headersKV)/2)
	for i := 0; i < len(headersKV); i += 2 {
		headers = append(headers, sarama.RecordHeader{
			Key:   []byte(headersKV[i]),
			Value: []byte(headersKV[i+1]),
		})
	}

	return &sarama.ProducerMessage{
		Topic:   topic,
		Key:     sarama.StringEncoder(key),
		Value:   sarama.ByteEncoder(message),
		Headers: headers,
	}, nil
}

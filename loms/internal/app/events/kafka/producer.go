package kafka

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"route256/loms/internal/app/models"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Producer struct {
	producer  sarama.AsyncProducer
	successes chan models.OutboxMessage
	fails     chan models.OutboxFailedMessage
}

func (p *Producer) ProduceMessage(ctx context.Context, message models.OutboxMessage) error {
	msg, err := p.BuildMessage(message.Destination, message.Key, message.Data)
	if err != nil {
		return fmt.Errorf("failed to build kafka message: %w", err)
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case p.producer.Input() <- msg:
		return nil
	}
}

func (p *Producer) Successes() <-chan models.OutboxMessage {
	return p.successes
}

func (p *Producer) Fails() <-chan models.OutboxFailedMessage {
	return p.fails
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

	producer := &Producer{
		producer:  asyncProducer,
		successes: make(chan models.OutboxMessage),
		fails:     make(chan models.OutboxFailedMessage),
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range asyncProducer.Errors() {
			key, err := msg.Msg.Key.Encode()
			if err != nil {
				log.Err(err).Msg("failed to encode kafka message key")
				continue
			}

			data, err := msg.Msg.Value.Encode()
			if err != nil {
				log.Err(err).Msg("failed to encode kafka message value")
				continue
			}

			outboxFailedMessage := models.OutboxFailedMessage{
				Message: models.OutboxMessage{
					Key:         string(key),
					Destination: msg.Msg.Topic,
					Data:        data,
				},
				Error: msg.Err,
			}

			log.Info().Msg(fmt.Sprintf("Async fail with key %s", outboxFailedMessage.Message.Key))
			producer.fails <- outboxFailedMessage
		}
		close(producer.fails)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range asyncProducer.Successes() {
			key, err := msg.Key.Encode()
			if err != nil {
				log.Err(err).Msg("failed to encode kafka message key")
				continue
			}

			data, err := msg.Value.Encode()
			if err != nil {
				log.Err(err).Msg("failed to encode kafka message value")
				continue
			}

			outboxMessage := models.OutboxMessage{
				Key:         string(key),
				Destination: msg.Topic,
				Data:        data,
			}

			log.Info().Msg(fmt.Sprintf("Async success with key %s", outboxMessage.Key))
			orderId, err := strconv.ParseInt(outboxMessage.Key, 10, 64)
			if err != nil {
				log.Err(err).Msg("failed to decode message key")
				continue
			}
			outboxMessage.Id = orderId
			producer.successes <- outboxMessage
		}
		close(producer.successes)
	}()

	return producer, nil
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

package kafka

import (
	"context"
	"fmt"
	"route256/loms/internal/app/services"
	"sync"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Producer struct {
	producer  sarama.AsyncProducer
	successes chan services.OutboxMessage
	fails     chan services.OutboxFailedMessage
}

func (p *Producer) ProduceMessage(ctx context.Context, message services.OutboxMessage) error {
	msg, err := p.BuildMessage(message.GetTopic(), message.GetKey(), message.GetData())
	if err != nil {
		return fmt.Errorf("failed to build kafka message: %w", err)
	}
	p.producer.Input() <- msg
	return nil
}

func (p *Producer) Successes() <-chan services.OutboxMessage {
	return p.successes
}

func (p *Producer) Fails() <-chan services.OutboxFailedMessage {
	return p.fails
}

type Message struct {
	Key   string
	Topic string
	Data  []byte
}

func (k Message) GetKey() string {
	return k.Key
}

func (k Message) GetData() []byte {
	return k.Data
}

func (k Message) GetTopic() string {
	return k.Topic
}

type FailedMessage struct {
	Message Message
	Error   error
}

func (f FailedMessage) GetMessage() services.OutboxMessage {
	return f.Message
}

func (f FailedMessage) GetError() error {
	return f.Error
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
		producer: asyncProducer,
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

			outboxFailedMessage := FailedMessage{
				Message: Message{
					Key:   string(key),
					Topic: msg.Msg.Topic,
					Data:  data,
				},
				Error: msg.Err,
			}

			log.Info().Msg(fmt.Sprintf("Async fail with key %s", outboxFailedMessage.Message.Key))
			producer.fails <- outboxFailedMessage
		}
		close(producer.successes)
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

			outboxMessage := Message{
				Key:   string(key),
				Topic: msg.Topic,
				Data:  data,
			}

			log.Info().Msg(fmt.Sprintf("Async success with key %s", outboxMessage.Key))
			producer.successes <- outboxMessage
		}
		close(producer.successes)
	}()

	return producer, nil
}

func (p *Producer) BuildOutboxMessage(ctx context.Context, key string, data []byte, topic string) (services.OutboxMessage, error) {
	return Message{
		Key:   key,
		Topic: topic,
		Data:  data,
	}, nil
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

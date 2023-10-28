package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"route256/loms/internal/app/models"
	"strconv"
	"sync"
	"time"
)

var OrderStatusChangedTopicName = "order-status-changed"

type OutboxRepo interface {
	SaveMessage(ctx context.Context, message OutboxMessage) error
	ClearLocks(outboxId string) error
	LockUnsentMessages(ctx context.Context, outboxId string) error
	GetLockedUnsentMessages(ctx context.Context, outboxId string) ([]OutboxMessage, error)
	SetMessageSent(ctx context.Context, message OutboxMessage) error
	SetMessageFailed(ctx context.Context, message OutboxMessage, err error) error
	GetFailedMessages(ctx context.Context, outboxId string) ([]OutboxFailedMessage, error)
}

type EventProducer interface {
	BuildOutboxMessage(ctx context.Context, key string, data []byte, topic string) (OutboxMessage, error)
	ProduceMessage(ctx context.Context, message OutboxMessage) error
	Successes() <-chan OutboxMessage
	Fails() <-chan OutboxFailedMessage
}

type Outbox struct {
	repo     OutboxRepo
	producer EventProducer
	id       string
}

type OutboxMessage interface {
	GetKey() string
	GetData() []byte
	GetTopic() string
}

type OutboxFailedMessage interface {
	GetMessage() OutboxMessage
	GetError() error
}

type orderStatusInfo struct {
	OrderId   int64              `json:"order_id"`
	NewStatus models.OrderStatus `json:"new_status"`
}

func (o *Outbox) OrderStatusChangedEventEmit(ctx context.Context, order models.Order) error {
	orderInfo := orderStatusInfo{
		OrderId:   order.Id,
		NewStatus: order.Status,
	}

	bytes, err := json.Marshal(orderInfo)
	if err != nil {
		log.Err(err).Msg("failed to marshal order info")
		return fmt.Errorf("failed to marshal order info: %w", err)
	}

	message, err := o.producer.BuildOutboxMessage(
		ctx,
		strconv.FormatInt(order.Id, 10),
		bytes,
		OrderStatusChangedTopicName,
	)
	if err != nil {
		log.Err(err).Msg("failed to build outbox message")
		return fmt.Errorf("failed to build outbox message: %w", err)
	}

	err = o.SaveMessage(ctx, message)
	if err != nil {
		log.Err(err).Msg("failed to save outbox message")
		return fmt.Errorf("failed to save outbox message: %w", err)
	}

	return nil
}

func (o *Outbox) SaveMessage(ctx context.Context, message OutboxMessage) error {
	err := o.repo.SaveMessage(ctx, message)
	if err != nil {
		return fmt.Errorf("error saving outbox message: %w", err)
	}
	return nil
}

func (o *Outbox) StartSendingMessages(ctx context.Context, wg *sync.WaitGroup, sendInterval time.Duration) {
	defer wg.Done()

	wg.Add(1)
	go o.ProcessSuccessfullySentMessages(ctx, wg)

	wg.Add(1)
	go o.ProcessFailedMessages(ctx, wg)

	ticker := time.NewTicker(sendInterval)
	for {
		select {
		case <-ticker.C:
			err := o.ProcessUnsentMessages(ctx)
			if err != nil {
				log.Error().Err(err).Msg("failed to process unsent outbox messages")
			}
		case <-ctx.Done():
			log.Info().Msg("stopping sending outbox messages...")
			ticker.Stop()
			log.Info().Msg("clearing outbox locks...")
			err := o.repo.ClearLocks(o.id)
			if err != nil {
				log.Error().Err(err).Msg("failed to clear locks for outbox")
			}
			log.Info().Msg("stopped producing messages")
			return
		}
	}
}

func (o *Outbox) ProcessUnsentMessages(ctx context.Context) error {
	err := o.repo.LockUnsentMessages(ctx, o.id)
	if err != nil {
		return fmt.Errorf("failed to lock unsent messages: %w", err)
	}

	unsentLockedMessages, err := o.repo.GetLockedUnsentMessages(ctx, o.id)
	if err != nil {
		return fmt.Errorf("failed to fetch locked unsent messages: %w", err)
	}

	for _, message := range unsentLockedMessages {
		err = o.producer.ProduceMessage(ctx, message)
		if err != nil {
			return fmt.Errorf("failed to produce messages: %w", err)
		}
	}

	return nil
}

func (o *Outbox) ProcessSuccessfullySentMessages(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Info().Msg("started processing successfully sent messages")
	for {
		select {
		case message := <-o.producer.Successes():
			err := o.repo.SetMessageSent(ctx, message)
			if err != nil {
				log.Error().Err(err).Msg("failed to set outbox message as sent")
				continue
			}

		case <-ctx.Done():
			log.Info().Msg("stopped processing successfully sent messages")
			return
		}
	}
}

func (o *Outbox) ProcessFailedMessages(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Info().Msg("started processing unsent messages")
	for {
		select {
		case failedMessage := <-o.producer.Fails():
			err := o.repo.SetMessageFailed(ctx, failedMessage.GetMessage(), failedMessage.GetError())
			if err != nil {
				log.Error().Err(err).Msg("failed to set outbox message as sent")
			}
		case <-ctx.Done():
			log.Info().Msg("stopped processing successfully sent messages")
			return
		}
	}
}

func (o *Outbox) StartRetryingFailedMessages(ctx context.Context, wg *sync.WaitGroup, retryInterval time.Duration) {
	defer wg.Done()
	log.Info().Msg("starting retrying sending failed outbox messages")

	ticker := time.NewTicker(retryInterval)
	for {
		select {
		case <-ticker.C:
			unsentFailedMessages, err := o.repo.GetFailedMessages(ctx, o.id)
			if err != nil {
				log.Error().Err(err).Msg("failed to fetch failed unsent messages")
				continue
			}

			for _, message := range unsentFailedMessages {
				err = o.producer.ProduceMessage(ctx, message.GetMessage())
				if err != nil {
					log.Error().Err(err).Msg("failed to produce messages")
					continue
				}
			}

		case <-ctx.Done():
			log.Info().Msg("stopped retrying sending failed outbox messages")
			ticker.Stop()
			return
		}
	}
}

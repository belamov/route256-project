package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"route256/loms/internal/app/models"
)

type Outbox struct {
	repo     MessagesProvider
	producer MessagesProducer
	id       string
}

func NewOutbox(outboxId string, producer MessagesProducer, provider MessagesProvider) *Outbox {
	return &Outbox{
		repo:     provider,
		producer: producer,
		id:       outboxId,
	}
}

var OrderStatusChangedTopicName = "order-status-changed"

type MessagesProvider interface {
	SaveMessage(ctx context.Context, message models.OutboxMessage) (models.OutboxMessage, error)
	ClearLocks(ctx context.Context, outboxId string) error
	LockUnsentMessages(ctx context.Context, outboxId string) error
	GetLockedUnsentMessages(ctx context.Context, outboxId string) ([]models.OutboxMessage, error)
	SetMessageSent(ctx context.Context, message models.OutboxMessage) error
	SetMessageFailed(ctx context.Context, message models.OutboxMessage, err error) error
	GetFailedMessages(ctx context.Context, outboxId string) ([]models.OutboxFailedMessage, error)
}

type MessagesProducer interface {
	ProduceMessage(ctx context.Context, message models.OutboxMessage) error
	Successes() <-chan models.OutboxMessage
	Fails() <-chan models.OutboxFailedMessage
}

func (o *Outbox) OrderStatusChangedEventEmit(ctx context.Context, order models.Order) error {
	orderInfo := models.OrderStatusInfo{
		OrderId:   order.Id,
		NewStatus: order.Status,
	}

	bytes, err := json.Marshal(orderInfo)
	if err != nil {
		log.Err(err).Msg("failed to marshal order info")
		return fmt.Errorf("failed to marshal order info: %w", err)
	}

	message := models.OutboxMessage{
		Key:         strconv.FormatInt(order.Id, 10),
		Destination: OrderStatusChangedTopicName,
		Data:        bytes,
	}

	_, err = o.SaveMessage(ctx, message)
	if err != nil {
		log.Err(err).Msg("failed to save outbox message")
		return fmt.Errorf("failed to save outbox message: %w", err)
	}

	return nil
}

func (o *Outbox) SaveMessage(ctx context.Context, message models.OutboxMessage) (models.OutboxMessage, error) {
	message, err := o.repo.SaveMessage(ctx, message)
	if err != nil {
		return message, fmt.Errorf("error saving outbox message: %w", err)
	}
	return message, nil
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
			err := o.repo.ClearLocks(context.Background(), o.id)
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
			log.Info().Any("msg", message).Msg("success message received")
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
			err := o.repo.SetMessageFailed(ctx, failedMessage.Message, failedMessage.Error)
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
				err = o.producer.ProduceMessage(ctx, message.Message)
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

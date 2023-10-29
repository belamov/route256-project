package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"route256/loms/internal/app/models"
	"route256/loms/internal/app/storage/repositories/queries"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OutboxPgRepository struct {
	dbPool       *pgxpool.Pool
	transactions pgTransactions
}

func NewOutboxPgRepository(dbPool *pgxpool.Pool) *OutboxPgRepository {
	return &OutboxPgRepository{
		dbPool:       dbPool,
		transactions: pgTransactions{},
	}
}

func (o *OutboxPgRepository) SaveMessage(ctx context.Context, message models.OutboxMessage) error {
	q := o.transactions.getQueriesFromContext(ctx, o.dbPool)

	params := queries.SaveOutboxMessageParams{
		Key:         message.Key,
		Destination: message.Destination,
		Data:        message.Data,
	}
	err := q.SaveOutboxMessage(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to save outbox message: %w", err)
	}

	return nil
}

func (o *OutboxPgRepository) ClearLocks(ctx context.Context, outboxId string) error {
	q := o.transactions.getQueriesFromContext(ctx, o.dbPool)

	return q.UnlockUnsentMessages(ctx, o.pgString(outboxId))
}

func (o *OutboxPgRepository) LockUnsentMessages(ctx context.Context, outboxId string) error {
	q := o.transactions.getQueriesFromContext(ctx, o.dbPool)

	return q.LockUnsentMessages(ctx, o.pgString(outboxId))
}

func (o *OutboxPgRepository) GetLockedUnsentMessages(ctx context.Context, outboxId string) ([]models.OutboxMessage, error) {
	q := o.transactions.getQueriesFromContext(ctx, o.dbPool)

	lockedMessages, err := q.GetLockedUnsentMessage(ctx, o.pgString(outboxId))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch locked messages: %w", err)
	}

	outboxMessages := make([]models.OutboxMessage, 0, len(lockedMessages))

	for _, lockedMessage := range lockedMessages {
		outboxMessages = append(outboxMessages, o.mapMessageToModel(lockedMessage))
	}

	return outboxMessages, nil
}

func (o *OutboxPgRepository) SetMessageSent(ctx context.Context, message models.OutboxMessage) error {
	q := o.transactions.getQueriesFromContext(ctx, o.dbPool)

	return q.SetMessageSent(ctx, message.Id)
}

func (o *OutboxPgRepository) SetMessageFailed(ctx context.Context, message models.OutboxMessage, err error) error {
	q := o.transactions.getQueriesFromContext(ctx, o.dbPool)

	params := queries.SetMessageFailedParams{
		ErrorMessage: o.pgString(err.Error()),
		ID:           message.Id,
	}
	return q.SetMessageFailed(ctx, params)
}

func (o *OutboxPgRepository) GetFailedMessages(ctx context.Context, outboxId string) ([]models.OutboxFailedMessage, error) {
	q := o.transactions.getQueriesFromContext(ctx, o.dbPool)

	failedMessages, err := q.GetFailedMessage(ctx, o.pgString(outboxId))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch failed messages: %w", err)
	}

	outboxFailedMessages := make([]models.OutboxFailedMessage, 0, len(failedMessages))
	for _, failedMessage := range failedMessages {
		outboxFailedMessages = append(outboxFailedMessages, models.OutboxFailedMessage{
			Message: o.mapMessageToModel(failedMessage),
			Error:   errors.New(failedMessage.ErrorMessage.String),
		})
	}

	return outboxFailedMessages, nil
}

func (o *OutboxPgRepository) pgString(value string) pgtype.Text {
	return pgtype.Text{
		String: value,
		Valid:  true,
	}
}

func (o *OutboxPgRepository) mapMessageToModel(messageFromTable queries.Outbox) models.OutboxMessage {
	return models.OutboxMessage{
		Id:          messageFromTable.ID,
		Key:         messageFromTable.Key,
		Destination: messageFromTable.Destination,
		Data:        messageFromTable.Data,
	}
}

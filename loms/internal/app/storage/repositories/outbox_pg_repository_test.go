package repositories

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"

	"route256/loms/internal/app/models"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"route256/loms/internal/app"
)

type OutboxPgRepositoryTestSuite struct {
	suite.Suite
	repo   *OutboxPgRepository
	dbPool *pgxpool.Pool
}

func TestOutboxPgRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(OutboxPgRepositoryTestSuite))
}

func (t *OutboxPgRepositoryTestSuite) SetupSuite() {
	config := app.BuildConfig()

	dbPool, err := InitPostgresDbConnection(config)
	require.NoError(t.T(), err)

	repo := NewOutboxPgRepository(dbPool)
	require.NoError(t.T(), err)
	t.repo = repo
	t.dbPool = dbPool
}

func (t *OutboxPgRepositoryTestSuite) TearDownSuite() {
	t.dbPool.Close()
}

func (t *OutboxPgRepositoryTestSuite) TestRepository() {
	ctx := context.Background()

	_, err := t.repo.dbPool.Exec(ctx, "truncate table outbox")
	require.NoError(t.T(), err)

	outboxId := strconv.Itoa(rand.Int())
	err = t.repo.ClearLocks(ctx, outboxId)
	assert.NoError(t.T(), err)

	message := models.OutboxMessage{
		Key:         "key",
		Destination: "topic",
		Data:        []byte("some data"),
	}

	message, err = t.repo.SaveMessage(ctx, message)
	assert.NoError(t.T(), err)

	sentMessage := models.OutboxMessage{
		Key:         "key",
		Destination: "topic",
		Data:        []byte("some data"),
	}
	savedSendMessage, err := t.repo.SaveMessage(ctx, sentMessage)
	assert.NoError(t.T(), err)

	err = t.repo.SetMessageSent(ctx, savedSendMessage)
	assert.NoError(t.T(), err)

	err = t.repo.LockUnsentMessages(ctx, outboxId)
	assert.NoError(t.T(), err)

	lockedMessages, err := t.repo.GetLockedUnsentMessages(ctx, outboxId)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), lockedMessages, 1)
	assert.NotEqual(t.T(), savedSendMessage.Id, lockedMessages[0].Id)
	assert.Equal(t.T(), message.Id, lockedMessages[0].Id)

	anotherMessage := models.OutboxMessage{
		Key:         "another",
		Destination: "topic",
		Data:        []byte("some another data"),
	}
	anotherMessage, err = t.repo.SaveMessage(ctx, anotherMessage)
	assert.NoError(t.T(), err)

	err = t.repo.LockUnsentMessages(ctx, outboxId)
	assert.NoError(t.T(), err)

	err = t.repo.SetMessageFailed(ctx, anotherMessage, errors.New("some error"))
	assert.NoError(t.T(), err)

	failedMessages, err := t.repo.GetFailedMessages(ctx, outboxId)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), failedMessages, 1)
	assert.Equal(t.T(), anotherMessage.Id, failedMessages[0].Message.Id)

	err = t.repo.ClearLocks(ctx, outboxId)
	assert.NoError(t.T(), err)

	lockedMessages, err = t.repo.GetLockedUnsentMessages(ctx, outboxId)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), lockedMessages, 0)
}

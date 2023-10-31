package services

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"route256/loms/internal/app/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type OutboxTestSuite struct {
	suite.Suite
	mockCtrl             *gomock.Controller
	mockMessagesProvider *MockMessagesProvider
	mockMessagesProducer *MockMessagesProducer
	outbox               *Outbox
}

var outboxId = "outboxId"

func (ts *OutboxTestSuite) SetupSuite() {
	ts.mockCtrl = gomock.NewController(Reporter{ts.T()})
	ts.mockMessagesProducer = NewMockMessagesProducer(ts.mockCtrl)
	ts.mockMessagesProvider = NewMockMessagesProvider(ts.mockCtrl)
	ts.outbox = NewOutbox(outboxId, ts.mockMessagesProducer, ts.mockMessagesProvider)
}

func TestOutboxTestSuite(t *testing.T) {
	suite.Run(t, new(OutboxTestSuite))
}

func (ts *OutboxTestSuite) TestOrderStatusChangedEventEmit() {
	ctx := context.Background()

	order := models.Order{
		Id:     1,
		Status: models.OrderStatusAwaitingPayment,
	}

	orderInfo := models.OrderStatusInfo{
		OrderId:   order.Id,
		NewStatus: order.Status,
	}
	bytes, err := json.Marshal(orderInfo)
	require.NoError(ts.T(), err)

	expectedMessage := models.OutboxMessage{
		Key:         strconv.FormatInt(order.Id, 10),
		Destination: OrderStatusChangedTopicName,
		Data:        bytes,
	}
	ts.mockMessagesProvider.EXPECT().SaveMessage(gomock.Any(), expectedMessage)

	err = ts.outbox.OrderStatusChangedEventEmit(ctx, order)
	assert.NoError(ts.T(), err)
}

func (ts *OutboxTestSuite) TestProducingMessages() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	ts.mockMessagesProvider.EXPECT().LockUnsentMessages(gomock.Any(), outboxId).Return(nil).AnyTimes()
	ts.mockMessagesProvider.EXPECT().GetLockedUnsentMessages(gomock.Any(), outboxId).Return(nil, nil).AnyTimes()
	ts.mockMessagesProvider.EXPECT().SetMessageSent(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	ts.mockMessagesProvider.EXPECT().ClearLocks(gomock.Any(), outboxId).Return(nil)

	ts.mockMessagesProducer.EXPECT().Fails().Return(make(chan models.OutboxFailedMessage))
	ts.mockMessagesProducer.EXPECT().Successes().Return(make(chan models.OutboxMessage))

	wg.Add(1)
	go ts.outbox.StartSendingMessages(ctx, wg, time.Millisecond*5)

	time.Sleep(time.Millisecond * 10)

	cancel()
	wg.Wait()
}

func (ts *OutboxTestSuite) TestRetryingMessages() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	failedMessage := models.OutboxFailedMessage{
		Message: models.OutboxMessage{
			Id:          2,
			Key:         "key",
			Destination: "topic",
			Data:        []byte("data"),
		},
		Error: errors.New("some error"),
	}
	failedMessages := []models.OutboxFailedMessage{failedMessage}

	ts.mockMessagesProvider.EXPECT().GetFailedMessages(gomock.Any(), outboxId).Return(failedMessages, nil).AnyTimes()
	ts.mockMessagesProducer.EXPECT().ProduceMessage(gomock.Any(), failedMessage.Message).AnyTimes()

	wg.Add(1)
	go ts.outbox.StartRetryingFailedMessages(ctx, wg, time.Millisecond*5)

	time.Sleep(time.Millisecond * 10)

	cancel()
	wg.Wait()
}

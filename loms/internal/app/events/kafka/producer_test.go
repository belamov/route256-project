package kafka

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"route256/loms/internal/app/models"
)

func TestProducer_OrderStatusChangedEventEmit(t *testing.T) {
	t.Skipf("integration test requires kafka up and running")

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	brokers := []string{"127.0.0.1:9091"} // from docker-compose
	wg.Add(1)

	producer, err := NewKafkaEventProducer(ctx, wg, brokers)
	require.NoError(t, err)

	time.Sleep(time.Second)
	changedOrder := models.Order{
		Id:     10,
		Status: models.OrderStatusPayed,
	}
	producer.OrderStatusChangedEventEmit(ctx, changedOrder)

	cancel()
	wg.Wait()
}

package kafka

import (
	"context"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"route256/loms/internal/app"
	"route256/loms/internal/app/services"
	"route256/loms/internal/app/storage/repositories"

	"github.com/stretchr/testify/require"
	"route256/loms/internal/app/models"
)

func TestProducer_OrderStatusChangedEventEmit(t *testing.T) {
	t.Skipf("integration test requires kafka and db up and running")
	config := app.BuildConfig()
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	brokers := []string{"127.0.0.1:9091"} // from docker-compose
	wg.Add(1)

	producer, err := NewKafkaEventProducer(ctx, wg, brokers)
	require.NoError(t, err)

	pgPool, err := repositories.InitPostgresDbConnection(ctx, config)
	require.NoError(t, err)

	outboxRepo := repositories.NewOutboxPgRepository(pgPool)

	orderEventProvider := services.NewOutbox(strconv.Itoa(rand.Int()), producer, outboxRepo)

	time.Sleep(time.Second)
	changedOrder := models.Order{
		Id:     10,
		Status: models.OrderStatusPayed,
	}
	err = orderEventProvider.OrderStatusChangedEventEmit(ctx, changedOrder)
	require.NoError(t, err)

	cancel()
	wg.Wait()
	pgPool.Close()
}

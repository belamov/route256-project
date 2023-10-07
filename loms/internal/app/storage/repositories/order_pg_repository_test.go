package repositories

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"route256/loms/internal/app"
	"route256/loms/internal/app/models"
	"route256/loms/internal/app/storage"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type OrderPgRepositoryTestSuite struct {
	suite.Suite
	repo   *OrderPgRepository
	cancel context.CancelFunc
	wg     *sync.WaitGroup
}

func TestOrderPgRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(OrderPgRepositoryTestSuite))
}

func (t *OrderPgRepositoryTestSuite) SetupSuite() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	config := app.BuildConfig()

	wg.Add(1)
	dbPool, err := InitPostgresDbConnection(ctx, wg, config)
	require.NoError(t.T(), err)

	repo := NewOrderPgRepository(dbPool)
	require.NoError(t.T(), err)
	t.repo = repo
	t.cancel = cancel
	t.wg = wg
}

func (t *OrderPgRepositoryTestSuite) TearDownSuite() {
	t.cancel()
	t.wg.Wait()
}

func (t *OrderPgRepositoryTestSuite) TestRepository() {
	ctx := context.Background()

	sku := rand.Uint32()

	order := models.Order{
		CreatedAt: time.Now(),
		Items: []models.OrderItem{{
			Name:  "sku",
			User:  1,
			Sku:   sku,
			Price: 10,
			Count: 1,
		}},
		UserId: 1,
		Status: models.OrderStatusNew,
	}

	createdOrder, err := t.repo.Create(ctx, order.UserId, order.Status, order.Items)
	assert.NoError(t.T(), err)
	order.Id = createdOrder.Id
	order.CreatedAt = createdOrder.CreatedAt
	assert.Equal(t.T(), order, createdOrder)

	updatedOrder, err := t.repo.SetStatus(ctx, createdOrder, models.OrderStatusCancelled)
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), models.OrderStatusCancelled, updatedOrder.Status)

	fetchedOrder, err := t.repo.GetOrderByOrderId(ctx, updatedOrder.Id)
	assert.NoError(t.T(), err)
	updatedOrder.CreatedAt = fetchedOrder.CreatedAt
	assert.Equal(t.T(), updatedOrder, fetchedOrder)

	_, err = t.repo.dbPool.Exec(ctx, "truncate table orders cascade")
	require.NoError(t.T(), err)

	_, err = t.repo.GetOrderByOrderId(ctx, 1)
	assert.ErrorIs(t.T(), err, storage.ErrOrderNotFound)

	// test fetching expired orders
	expiredOrder, err := t.repo.Create(ctx, order.UserId, models.OrderStatusAwaitingPayment, order.Items)
	assert.NoError(t.T(), err)

	expiredOrdersIds, err := t.repo.GetExpiredOrdersWithStatus(
		ctx,
		time.Now().Add(time.Hour*100),
		models.OrderStatusAwaitingPayment,
	)
	assert.NoError(t.T(), err)
	require.Len(t.T(), expiredOrdersIds, 1)
	assert.Equal(t.T(), expiredOrder.Id, expiredOrdersIds[0])
}

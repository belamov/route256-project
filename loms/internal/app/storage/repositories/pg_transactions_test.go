package repositories

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"route256/loms/internal/app"
	"route256/loms/internal/app/models"
	"route256/loms/internal/app/storage"
)

type PgTransactionsTestSuite struct {
	suite.Suite
	transactor *PgTransactor
	stocksRepo *StocksPgRepository
	ordersRepo *OrderPgRepository
	cancel     context.CancelFunc
	wg         *sync.WaitGroup
}

func TestPgTransactionsTestSuite(t *testing.T) {
	suite.Run(t, new(PgTransactionsTestSuite))
}

func (t *PgTransactionsTestSuite) SetupSuite() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	config := app.BuildConfig()

	wg.Add(1)
	dbPool, err := InitPostgresDbConnection(ctx, wg, config)
	require.NoError(t.T(), err)

	t.ordersRepo = NewOrderPgRepository(dbPool)
	t.stocksRepo = NewStocksPgRepository(dbPool)
	t.transactor = NewPgTransactor(dbPool)
	t.cancel = cancel
	t.wg = wg
}

func (t *PgTransactionsTestSuite) TearDownSuite() {
	t.cancel()
	t.wg.Wait()
}

func (t *PgTransactionsTestSuite) TestTransactions() {
	ctx := context.Background()

	var notFoundSku uint32 = 1

	order := models.Order{
		CreatedAt: time.Now(),
		Items: []models.OrderItem{{
			Name:  "sku",
			User:  1,
			Sku:   notFoundSku,
			Price: 10,
			Count: 1,
		}},
		UserId: 1,
		Status: models.OrderStatusNew,
	}

	createdOrder, err := t.ordersRepo.Create(ctx, order.UserId, order.Status, order.Items)
	assert.NoError(t.T(), err)

	err = t.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		updatedOrder, err := t.ordersRepo.SetStatus(ctx, createdOrder, models.OrderStatusAwaitingPayment)
		assert.NoError(t.T(), err)

		_, err = t.stocksRepo.dbPool.Exec(ctx, "delete from stocks where sku = $1", int(notFoundSku))
		assert.NoError(t.T(), err)

		err = t.stocksRepo.Reserve(ctx, updatedOrder)
		assert.ErrorIs(t.T(), err, storage.ErrInsufficientStocks)
		return err
	})
	assert.Error(t.T(), err)

	fetchedOrder, err := t.ordersRepo.GetOrderByOrderId(ctx, createdOrder.Id)
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), fetchedOrder.Status, models.OrderStatusNew)
}

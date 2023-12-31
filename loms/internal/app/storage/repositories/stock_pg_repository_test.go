package repositories

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"route256/loms/internal/app/storage"

	"github.com/stretchr/testify/assert"
	"route256/loms/internal/app"
	"route256/loms/internal/app/models"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StockPgRepositoryTestSuite struct {
	suite.Suite
	repo   *StocksPgRepository
	dbPool *pgxpool.Pool
}

func TestStockPgRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(StockPgRepositoryTestSuite))
}

func (t *StockPgRepositoryTestSuite) SetupSuite() {
	config := app.BuildConfig()

	dbPool, err := InitPostgresDbConnection(context.Background(), config)
	require.NoError(t.T(), err)

	t.dbPool = dbPool

	repo := NewStocksPgRepository(dbPool)
	require.NoError(t.T(), err)
	t.repo = repo
}

func (t *StockPgRepositoryTestSuite) TearDownSuite() {
	t.dbPool.Close()
}

func (t *StockPgRepositoryTestSuite) TestRepository() {
	ctx := context.Background()
	t.clearStocks(ctx)

	sku := rand.Uint32()
	var initialStockCount uint64 = 10

	t.addStocksForSku(ctx, sku, initialStockCount)
	t.addStocksForSku(ctx, sku+1, initialStockCount)

	count, err := t.repo.GetBySku(ctx, sku)
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), initialStockCount, count)

	order := models.Order{
		CreatedAt: time.Now(),
		Items: []models.OrderItem{{
			Name:  "sku",
			User:  1,
			Sku:   sku,
			Price: 10,
			Count: 1,
		}},
		Id:     1,
		UserId: 1,
		Status: models.OrderStatusNew,
	}

	err = t.repo.Reserve(ctx, order)
	assert.NoError(t.T(), err)

	count, err = t.repo.GetBySku(ctx, sku)
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), initialStockCount-1, count)

	err = t.repo.ReserveCancel(ctx, order)
	assert.NoError(t.T(), err)

	count, err = t.repo.GetBySku(ctx, sku)
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), initialStockCount, count)

	err = t.repo.ReserveCancel(ctx, order)
	assert.NoError(t.T(), err)

	count, err = t.repo.GetBySku(ctx, sku)
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), initialStockCount+1, count)
}

func (t *StockPgRepositoryTestSuite) TestRepositoryInsufficientStocks() {
	ctx := context.Background()

	t.clearStocks(ctx)

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
		Id:     1,
		UserId: 1,
		Status: models.OrderStatusNew,
	}

	err := t.repo.Reserve(ctx, order)
	assert.ErrorIs(t.T(), err, storage.ErrInsufficientStocks)

	t.addStocksForSku(ctx, sku, 1)

	err = t.repo.Reserve(ctx, order)
	assert.NoError(t.T(), err)

	err = t.repo.Reserve(ctx, order)
	assert.ErrorIs(t.T(), err, storage.ErrInsufficientStocks)
}

func (t *StockPgRepositoryTestSuite) clearStocks(ctx context.Context) {
	_, err := t.repo.dbPool.Exec(ctx, "truncate table stocks")
	require.NoError(t.T(), err)
}

func (t *StockPgRepositoryTestSuite) addStocksForSku(ctx context.Context, sku uint32, count uint64) {
	_, err := t.repo.dbPool.Exec(ctx, "insert into stocks (sku, count) values ($1, $2)", sku, count)
	require.NoError(t.T(), err)
}

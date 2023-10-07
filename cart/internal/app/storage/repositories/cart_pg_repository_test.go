package repositories

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"route256/cart/internal/app"
	"route256/cart/internal/app/models"
	"route256/cart/internal/app/services"
)

type PgRepositoryTestSuite struct {
	suite.Suite
	repo   services.CartProvider
	cancel context.CancelFunc
	wg     *sync.WaitGroup
}

func TestPgRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PgRepositoryTestSuite))
}

func initPostgresDbConnection(ctx context.Context, wg *sync.WaitGroup, config *app.Config) (*pgxpool.Pool, error) {
	databaseDSN := fmt.Sprintf(
		"postgresql://%s:%s@%s/%s",
		config.DbUser,
		config.DbPassword,
		config.DbHost,
		config.DbName,
	)
	dbPool, err := pgxpool.New(ctx, databaseDSN)
	if err != nil {
		return nil, err
	}
	log.Info().Msg("Connected to postgres")

	go func() {
		<-ctx.Done()
		log.Info().Msg("Closing order repository connections...")
		dbPool.Close()
		log.Info().Msg("Order repository connections closed")
		wg.Done()
	}()

	return dbPool, nil
}

func (t *PgRepositoryTestSuite) SetupSuite() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	config := app.BuildConfig()

	wg.Add(1)
	dbPool, err := initPostgresDbConnection(ctx, wg, config)
	require.NoError(t.T(), err)

	repo := NewCartRepository(dbPool)
	require.NoError(t.T(), err)
	t.repo = repo
	t.cancel = cancel
	t.wg = wg
}

func (t *PgRepositoryTestSuite) TearDownSuite() {
	t.cancel()
	t.wg.Wait()
}

func (t *PgRepositoryTestSuite) TestRepository() {
	maxId := rand.Int63()

	maxItem1 := models.CartItem{
		User:  maxId,
		Sku:   1,
		Count: 1,
	}

	maxItem2 := models.CartItem{
		User:  maxId,
		Sku:   2,
		Count: 5,
	}
	err := t.repo.SaveItem(context.Background(), maxItem1)
	assert.NoError(t.T(), err)
	err = t.repo.SaveItem(context.Background(), maxItem2)
	assert.NoError(t.T(), err)

	johnId := rand.Int63()
	johnItem1 := models.CartItem{
		User:  johnId,
		Sku:   1,
		Count: 1,
	}

	johnItem2 := models.CartItem{
		User:  johnId,
		Sku:   2,
		Count: 5,
	}
	err = t.repo.SaveItem(context.Background(), johnItem1)
	assert.NoError(t.T(), err)
	err = t.repo.SaveItem(context.Background(), johnItem2)
	assert.NoError(t.T(), err)

	fetchedItems, err := t.repo.GetItemsByUserId(context.Background(), maxId)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), fetchedItems, 2)

	fetchedItems, err = t.repo.GetItemsByUserId(context.Background(), maxId)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), fetchedItems, 2)

	err = t.repo.DeleteItem(context.Background(), maxItem1)
	assert.NoError(t.T(), err)

	fetchedItems, err = t.repo.GetItemsByUserId(context.Background(), maxId)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), fetchedItems, 1)

	err = t.repo.DeleteItemsByUserId(context.Background(), maxId)
	assert.NoError(t.T(), err)

	fetchedItems, err = t.repo.GetItemsByUserId(context.Background(), maxId)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), fetchedItems, 0)

	err = t.repo.DeleteItemsByUserId(context.Background(), johnId)
	assert.NoError(t.T(), err)

	fetchedItems, err = t.repo.GetItemsByUserId(context.Background(), johnId)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), fetchedItems, 0)
}

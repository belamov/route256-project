package cache

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"route256/cart/internal/app"
	"route256/cart/internal/app/models"
)

func TestRedisShardedCache(t *testing.T) {
	t.Skip("Skip integration test")

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	config := app.BuildConfig()

	wg.Add(1)
	cache := NewRedis(ctx, wg, config.RedisShards)

	for i := 0; i < 100; i++ {
		items := make([]models.CartItemWithInfo, 1)
		items[0] = models.CartItemWithInfo{
			Name:  strconv.Itoa(i),
			User:  int64(i),
			Sku:   uint32(i),
			Price: uint32(i),
			Count: uint64(i),
		}

		err := cache.SetCartItems(ctx, int64(i), items)
		assert.NoError(t, err)
	}

	for i := 0; i < 100; i++ {
		items, err := cache.GetCartItems(ctx, int64(i))
		assert.NoError(t, err)
		assert.Len(t, items, 1)
		assert.Equal(t, models.CartItemWithInfo{
			Name:  strconv.Itoa(i),
			User:  int64(i),
			Sku:   uint32(i),
			Price: uint32(i),
			Count: uint64(i),
		}, items[0])
	}

	notExistingItems, err := cache.GetCartItems(ctx, 100000000)
	assert.NoError(t, err)
	assert.Nil(t, notExistingItems)
	fmt.Println(notExistingItems, err)

	for i := 0; i < 100; i++ {
		err := cache.Invalidate(ctx, int64(i))
		assert.NoError(t, err)
	}

	for i := 0; i < 100; i++ {
		items, err := cache.GetCartItems(ctx, int64(i))
		assert.NoError(t, err)
		assert.Nil(t, items)
	}

	cancel()
	wg.Wait()
}

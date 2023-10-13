package product

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"route256/cart/internal/pkg/limiter"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProductWithRateLimiter(t *testing.T) {
	t.Skip("Skip real api call")

	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), "products_token", "testtoken"))

	wg := &sync.WaitGroup{}

	targetRps := 2
	goroutinesCount := 10
	var existingSku uint32 = 1148162
	var expectedPrice uint32 = 2931
	realProductUrl := "route256.pavl.uk:8082"

	rateLimiter := limiter.NewRateLimiter(targetRps)

	wg.Add(1)
	client, err := NewProductGrpcClient(ctx, wg, realProductUrl, rateLimiter)
	require.NoError(t, err)

	wgRequests := &sync.WaitGroup{}
	wgRequests.Add(goroutinesCount)

	for i := 0; i < goroutinesCount; i++ {
		go func() {
			defer wgRequests.Done()
			info, err := client.GetProduct(ctx, existingSku)
			fmt.Println(time.Now(), info, err)
			assert.NoError(t, err)
			assert.Equal(t, "Кулинар Гуров", info.Name)
			assert.Equal(t, expectedPrice, info.Price)
		}()
	}

	wgRequests.Wait()
	cancel()
	wg.Wait()
}

func TestGetProductRedisRateLimiter(t *testing.T) {
	t.Skip("Skip real api call")

	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), "products_token", "testtoken"))

	wg := &sync.WaitGroup{}

	targetRps := 2
	goroutinesCount := 10
	var existingSku uint32 = 1148162
	var expectedPrice uint32 = 2931
	realProductUrl := "route256.pavl.uk:8082"

	redisRateLimiter := limiter.NewRedisRateLimiter(ctx, wg, "redis:6379", targetRps)

	wg.Add(1)
	client, err := NewProductGrpcClient(ctx, wg, realProductUrl, redisRateLimiter)
	require.NoError(t, err)

	wgRequests := &sync.WaitGroup{}
	wgRequests.Add(goroutinesCount)

	for i := 0; i < goroutinesCount; i++ {
		go func() {
			defer wgRequests.Done()
			info, err := client.GetProduct(ctx, existingSku)
			fmt.Println(time.Now(), info, err)
			assert.NoError(t, err)
			assert.Equal(t, "Кулинар Гуров", info.Name)
			assert.Equal(t, expectedPrice, info.Price)
		}()
	}

	wgRequests.Wait()
	cancel()
	wg.Wait()
}

package product

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"route256/cart/internal/app/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProductWithRateLimiter(t *testing.T) {
	t.Skip("Skip real api call")

	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), "products_token", "testtoken"))

	wg := &sync.WaitGroup{}
	wg.Add(1)
	limiter := services.NewRateLimiter(2)
	client, err := NewProductGrpcClient(ctx, wg, "route256.pavl.uk:8082", limiter)
	require.NoError(t, err)

	wgRequests := &sync.WaitGroup{}
	wgRequests.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			defer wgRequests.Done()
			info, err := client.GetProduct(ctx, 1148162)
			fmt.Println(time.Now(), info, err)
			assert.NoError(t, err)
			assert.Equal(t, "Кулинар Гуров", info.Name)
			assert.Equal(t, uint32(2931), info.Price)
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
	wg.Add(1)

	limiter := services.NewRedisRateLimiter(ctx, wg, "redis:6379", 2)
	client, err := NewProductGrpcClient(ctx, wg, "route256.pavl.uk:8082", limiter)
	require.NoError(t, err)
	wgRequests := &sync.WaitGroup{}
	wgRequests.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			defer wgRequests.Done()
			info, err := client.GetProduct(ctx, 1148162)
			fmt.Println(time.Now(), info, err)
			assert.NoError(t, err)
			assert.Equal(t, "Кулинар Гуров", info.Name)
			assert.Equal(t, uint32(2931), info.Price)
		}()
	}

	wgRequests.Wait()
	cancel()
	wg.Wait()
}

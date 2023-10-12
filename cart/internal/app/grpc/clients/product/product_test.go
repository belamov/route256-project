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

func TestGetProduct(t *testing.T) {
	t.Skip("Skip real api call")

	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), "products_token", "testtoken"))

	wg := &sync.WaitGroup{}
	wg.Add(1)
	limiter := services.NewSinglePodLimiter(5)
	client, err := NewProductGrpcClient(ctx, wg, "route256.pavl.uk:8082", limiter)
	require.NoError(t, err)
	for i := 0; i < 20; i++ {
		info, err := client.GetProduct(ctx, 1148162)
		fmt.Println(time.Now(), info, err)
		assert.NoError(t, err)
		assert.Equal(t, "Кулинар Гуров", info.Name)
		assert.Equal(t, uint32(2931), info.Price)
	}

	cancel()
	wg.Wait()
}

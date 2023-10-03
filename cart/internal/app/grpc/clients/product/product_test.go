package product

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProduct(t *testing.T) {
	t.Skip("Skip real api call")

	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), "products_token", "testtoken"))

	wg := &sync.WaitGroup{}
	wg.Add(1)
	client, err := NewProductGrpcClient(ctx, wg, "route256.pavl.uk:8082")
	require.NoError(t, err)

	info, err := client.GetProduct(ctx, 1148162)
	assert.NoError(t, err)
	assert.Equal(t, "Кулинар Гуров", info.Name)
	assert.Equal(t, uint32(2931), info.Price)
	cancel()
	wg.Wait()
}

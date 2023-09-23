package clients

import (
	"context"
	"testing"

	"route256/cart/internal/app/domain/services"

	"github.com/stretchr/testify/assert"
)

func Test_productHttpClient_GetProduct(t *testing.T) {
	t.Skip("Skip real api call")

	client := NewProductHttpClient("http://route256.pavl.uk:8080/get_product")
	ctx := context.WithValue(context.Background(), "products_token", "testtoken")
	product, err := client.GetProduct(ctx, 773297411)
	assert.NoError(t, err)
	assert.Equal(t, product.Name, "Кроссовки Nike JORDAN")
	assert.Equal(t, product.Price, uint32(2202))

	_, err = client.GetProduct(ctx, 1)
	assert.ErrorIs(t, err, services.ErrSkuInvalid)
}

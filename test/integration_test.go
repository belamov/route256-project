package test

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/rand"
	pb "route256/test/cartbp"
	"sync"
	"testing"
)

func TestCart(t *testing.T) {
	// перед запуском в корне проекта выполнить
	// cp -l -r cart/internal/app/grpc/pb test/cartbp
	// make up-kafka
	// make run-all
	// make migrate
	// добавить стоки для скю 1148162 в таблицу loms.stocks
	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), "products_token", "testtoken"))
	wg := &sync.WaitGroup{}

	cartClient, err := NewCartGrpcClient(ctx, wg, "localhost:8083")
	require.NoError(t, err)

	for {

		userId := rand.Int63()
		var sku uint32 = 1148162
		addItemRequest := &pb.AddItemRequest{
			User: userId,
			Item: &pb.CartItemAddRequest{
				User:  userId,
				Sku:   sku,
				Count: 1,
			},
		}
		_, err = cartClient.AddItem(ctx, addItemRequest)
		assert.NoError(t, err)

		response, err := cartClient.List(ctx, &pb.ListRequest{User: userId})
		assert.NoError(t, err)
		assert.Len(t, response.Items, 1)
		assert.Equal(t, sku, response.Items[0].Sku)
		assert.Equal(t, sku, response.Items[0].Sku)
		assert.Equal(t, uint64(1), response.Items[0].Count)

		_, err = cartClient.DeleteItem(ctx, &pb.DeleteItemRequest{
			User: userId,
			Sku:  sku,
		})
		assert.NoError(t, err)

		response, err = cartClient.List(ctx, &pb.ListRequest{User: userId})
		assert.NoError(t, err)
		assert.Len(t, response.Items, 0)

		_, err = cartClient.AddItem(ctx, addItemRequest)
		assert.NoError(t, err)

		checkoutResponse, err := cartClient.Checkout(ctx, &pb.CheckoutRequest{User: userId})
		assert.NoError(t, err)
		assert.Greater(t, checkoutResponse.OrderID, int64(0))
	}

	cancel()
	wg.Wait()
}

func NewCartGrpcClient(ctx context.Context, wg *sync.WaitGroup, serviceUrl string) (pb.CartClient, error) {
	conn, err := grpc.DialContext(
		ctx,
		serviceUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	grpcClient := pb.NewCartClient(conn)

	go func() {
		<-ctx.Done()
		log.Info().Msg("Closing loms grpc client")
		err := conn.Close()
		if err != nil {
			log.Err(err).Msg("Couldn't close loms grpc connection")
			wg.Done()
			return
		}
		log.Info().Msg("Closed loms grpc client")
		wg.Done()
	}()

	log.Info().Msg("loms grpc client configured. connected to " + serviceUrl)

	return grpcClient, nil
}

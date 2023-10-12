package product

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"route256/cart/internal/app/grpc/clients/product/pb"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"route256/cart/internal/app/models"
	"route256/cart/internal/app/services"
)

type productGrpcClient struct {
	grpcClient pb.ProductServiceClient
	conn       *grpc.ClientConn
}

type Limiter interface {
	Wait(ctx context.Context) error
}

func UnaryClientInterceptor(limiter Limiter) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		err := limiter.Wait(ctx)
		if err != nil {
			return fmt.Errorf("cant wait for limiter to allow request: %w", err)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func NewProductGrpcClient(ctx context.Context, wg *sync.WaitGroup, serviceUrl string, limiter Limiter) (services.ProductService, error) {
	conn, err := grpc.DialContext(
		ctx,
		serviceUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(UnaryClientInterceptor(limiter)),
	)
	if err != nil {
		return nil, err
	}

	grpcClient := pb.NewProductServiceClient(conn)

	go func() {
		<-ctx.Done()
		log.Info().Msg("Closing product grpc client")
		err := conn.Close()
		if err != nil {
			log.Err(err).Msg("Couldn't close product grpc connection")
			wg.Done()
			return
		}
		log.Info().Msg("Closed product grpc client")
		wg.Done()
	}()

	log.Info().Msg("product grpc client configured. connected to " + serviceUrl)

	return &productGrpcClient{
		grpcClient: grpcClient,
		conn:       conn,
	}, nil
}

func (p *productGrpcClient) GetProduct(ctx context.Context, sku uint32) (models.CartItemInfo, error) {
	token, ok := ctx.Value("products_token").(string)
	if !ok {
		return models.CartItemInfo{}, errors.New("cant parse products_token from context")
	}

	request := &pb.GetProductRequest{
		Token: token,
		Sku:   sku,
	}

	response, err := p.grpcClient.GetProduct(ctx, request)
	if err != nil {
		return models.CartItemInfo{}, err
	}

	info := models.CartItemInfo{
		Name:  response.Name,
		Sku:   sku,
		Price: response.Price,
	}
	return info, nil
}

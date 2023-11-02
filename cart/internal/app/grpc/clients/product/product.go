package product

import (
	"context"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"sync"

	"route256/cart/internal/app/grpc/clients/product/pb"
	"route256/cart/internal/app/grpc/interceptors"
	"route256/cart/internal/app/models"
	"route256/cart/internal/app/services"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type productGrpcClient struct {
	grpcClient pb.ProductServiceClient
	conn       *grpc.ClientConn
}

func NewProductGrpcClient(
	ctx context.Context,
	wg *sync.WaitGroup,
	serviceUrl string,
	limiter interceptors.Limiter,
) (services.ProductService, error) {
	conn, err := grpc.DialContext(
		ctx,
		serviceUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptors.RateLimitClientInterceptor(limiter)),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
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
	request := &pb.GetProductRequest{
		Token: "testtoken",
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

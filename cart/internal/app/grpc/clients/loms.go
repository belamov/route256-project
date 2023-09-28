package clients

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"route256/cart/internal/app/models"
	"route256/cart/internal/app/services"
	lomspb "route256/loms/api/proto"
)

type lomsGrpcClient struct {
	grpcClient lomspb.LomsClient
	conn       *grpc.ClientConn
}

func NewLomsGrpcClient(ctx context.Context, wg *sync.WaitGroup, serviceUrl string) (services.LomsService, error) {
	conn, err := grpc.DialContext(
		ctx,
		serviceUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	grpcClient := lomspb.NewLomsClient(conn)

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

	return &lomsGrpcClient{
		grpcClient: grpcClient,
		conn:       conn,
	}, nil
}

func (l *lomsGrpcClient) GetStocksInfo(ctx context.Context, sku uint32) (uint64, error) {
	request := &lomspb.StockInfoRequest{Sku: sku}

	response, err := l.grpcClient.StockInfo(ctx, request)
	if err != nil {
		return 0, err
	}

	return response.Count, nil
}

func (l *lomsGrpcClient) CreateOrder(ctx context.Context, userId int64, items []models.CartItem) (int64, error) {
	request := &lomspb.OrderCreateRequest{
		User:  userId,
		Items: make([]*lomspb.OrderItemCreateRequest, 0, len(items)),
	}
	for _, item := range items {
		request.Items = append(request.Items, &lomspb.OrderItemCreateRequest{
			Sku:   item.Sku,
			Count: item.Count,
		})
	}

	response, err := l.grpcClient.OrderCreate(ctx, request)
	if err != nil {
		return 0, err
	}

	return response.OrderId, nil
}

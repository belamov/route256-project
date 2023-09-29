package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	"route256/loms/internal/app"
	grpcserver "route256/loms/internal/app/grpc/server"
	httpserver "route256/loms/internal/app/http/server"
	"route256/loms/internal/app/models"
	"route256/loms/internal/app/services"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type NullOrderProvider struct{}

func (n NullOrderProvider) GetOrdersIdsByCreatedAtAndStatus(ctx context.Context, createdAt time.Time, orderStatus models.OrderStatus) ([]int64, error) {
	return []int64{}, nil
}

func (n NullOrderProvider) Create(ctx context.Context, userId int64, statusNew models.OrderStatus, items []models.OrderItem) (models.Order, error) {
	// TODO implement me
	panic("implement me")
}

func (n NullOrderProvider) SetStatus(ctx context.Context, order models.Order, status models.OrderStatus) (models.Order, error) {
	// TODO implement me
	panic("implement me")
}

func (n NullOrderProvider) GetOrderByOrderId(ctx context.Context, orderId int64) (models.Order, error) {
	// TODO implement me
	panic("implement me")
}

func (n NullOrderProvider) CancelUnpaidOrders(ctx context.Context) error {
	return nil
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	config := app.BuildConfig()

	stocksProvider := services.NewMockStocksProvider(nil)
	lomsService := services.NewLomsService(NullOrderProvider{}, stocksProvider, config.AllowedOrderUnpaidTime)

	httpServer := httpserver.NewHTTPServer(config.HttpServerAddress, lomsService)
	grpcServer := grpcserver.NewGRPCServer(config.GrpcServerAddress, config.GrpcGatewayServerAddress, lomsService)

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	wg := &sync.WaitGroup{}

	wg.Add(4)

	go httpServer.Run(ctx, wg)
	go lomsService.RunCancelUnpaidOrders(ctx, wg, config.CancelUnpaidOrdersInterval)
	go grpcServer.Run(ctx, wg)
	go grpcServer.RunGateway(ctx, wg)

	wg.Wait()

	log.Info().Msg("goodbye")
}

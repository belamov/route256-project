package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	"route256/loms/internal/app/models"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"route256/loms/internal/app"
	mocks "route256/loms/internal/app/mocks"
	"route256/loms/internal/app/server"
	"route256/loms/internal/app/services"
)

type NullOrderProvider struct{}

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

	stocksProvider := mocks.NewMockStocksProvider(nil)
	lomsService := services.NewLomsService(NullOrderProvider{}, stocksProvider)

	config := app.BuildServerConfig()

	srv := server.NewHTTPServer(config.Address, lomsService)

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	wg := &sync.WaitGroup{}

	wg.Add(2)
	go srv.Run(ctx, wg)
	go lomsService.RunCancelUnpaidOrders(ctx, wg, time.Minute)
	<-ctx.Done()
	wg.Wait()
	log.Info().Msg("goodbye")
}

package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"os/signal"
	"route256/loms/internal/app"
	grpcserver "route256/loms/internal/app/grpc/server"
	httpserver "route256/loms/internal/app/http/server"
	"route256/loms/internal/app/services"
	"route256/loms/internal/app/storage/repositories/order"
	"route256/loms/internal/app/storage/repositories/stocks"
	"sync"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	config := app.BuildConfig()

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	dbPool, err := initPostgresDbConnection(ctx, wg, config)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize connection to postgres")
		return
	}

	stocksProvider := stocks.NewStocksPgRepository(dbPool)
	orderProvider := order.NewOrderRepository(dbPool)
	lomsService := services.NewLomsService(orderProvider, stocksProvider, config.AllowedOrderUnpaidTime)

	httpServer := httpserver.NewHTTPServer(config.HttpServerAddress, lomsService)
	grpcServer := grpcserver.NewGRPCServer(config.GrpcServerAddress, config.GrpcGatewayServerAddress, lomsService)

	wg.Add(4)

	go httpServer.Run(ctx, wg)
	go lomsService.RunCancelUnpaidOrders(ctx, wg, config.CancelUnpaidOrdersInterval)
	go grpcServer.Run(ctx, wg)
	go grpcServer.RunGateway(ctx, wg)

	wg.Wait()

	log.Info().Msg("goodbye")
}

func initPostgresDbConnection(ctx context.Context, wg *sync.WaitGroup, config *app.Config) (*pgxpool.Pool, error) {
	databaseDSN := fmt.Sprintf(
		"postgresql://%s:%s@%s/%s",
		config.DbUser,
		config.DbPassword,
		config.DbHost,
		config.DbName,
	)
	dbPool, err := pgxpool.New(ctx, databaseDSN)
	if err != nil {
		return nil, err
	}
	log.Info().Msg("Connected to postgres")

	go func() {
		<-ctx.Done()
		log.Info().Msg("Closing order repository connections...")
		dbPool.Close()
		log.Info().Msg("Order repository connections closed")
		wg.Done()
	}()

	return dbPool, nil
}

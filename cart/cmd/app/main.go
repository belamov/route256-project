package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"

	"route256/cart/internal/app/storage/repositories/cart"

	"route256/cart/internal/app/http/handlers"

	"route256/cart/internal/app"
	"route256/cart/internal/app/grpc/clients/loms"
	"route256/cart/internal/app/grpc/clients/product"
	grpcserver "route256/cart/internal/app/grpc/server"
	httpserver "route256/cart/internal/app/http/server"
	"route256/cart/internal/app/services"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	config := app.BuildConfig()

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	lomsService, err := loms.NewLomsGrpcClient(ctx, wg, config.LomsGrpcServiceUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("failed init grpc loms client")
		return
	}

	wg.Add(1)
	productService, err := product.NewProductGrpcClient(ctx, wg, config.ProductGrpcServiceUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("failed init grpc product client")
		return
	}

	wg.Add(1)
	dbPool, err := initPostgresDbConnection(ctx, wg, config)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize connection to postgres")
		return
	}

	cartProvider := cart.NewCartRepository(dbPool)

	cartService := services.NewCartService(productService, lomsService, cartProvider)

	httpServer := httpserver.NewHTTPServer(config.HttpServerAddress, handlers.NewRouter(cartService))
	grpcServer := grpcserver.NewGRPCServer(config.GrpcServerAddress, config.GrpcGatewayServerAddress, cartService)

	wg.Add(3)

	go httpServer.Run(ctx, wg)
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

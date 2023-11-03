package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"route256/cart/internal/pkg/tracer"

	"route256/cart/internal/app"
	"route256/cart/internal/app/grpc/clients/loms"
	"route256/cart/internal/app/grpc/clients/product"
	grpcserver "route256/cart/internal/app/grpc/server"
	"route256/cart/internal/app/http/handlers"
	httpserver "route256/cart/internal/app/http/server"
	"route256/cart/internal/app/services"
	"route256/cart/internal/app/storage/repositories"
	"route256/cart/internal/pkg/limiter"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Caller().Logger()
	zerolog.SetGlobalLevel(0)

	config := app.BuildConfig()

	zerolog.SetGlobalLevel(config.LogLevel)

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	_, err := tracer.InitTracer(ctx, wg, "localhost:4318", "", "cart")
	if err != nil {
		log.Fatal().Err(err).Msg("failed init tracer")
		return
	}

	wg.Add(1)
	lomsService, err := loms.NewLomsGrpcClient(ctx, wg, config.LomsGrpcServiceUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("failed init grpc loms client")
		return
	}

	wg.Add(1)
	redisRateLimiter := limiter.NewRedisRateLimiter(ctx, wg, config.RedisAddress, config.TargetRpsToProductService)

	wg.Add(1)
	productService, err := product.NewProductGrpcClient(ctx, wg, config.ProductGrpcServiceUrl, redisRateLimiter)
	if err != nil {
		log.Fatal().Err(err).Msg("failed init grpc product client")
		return
	}

	wg.Add(1)
	dbPool, err := repositories.InitPostgresDbConnection(ctx, wg, config)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize connection to postgres")
		return
	}

	cartProvider := repositories.NewCartRepository(dbPool)

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

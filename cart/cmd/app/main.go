package main

import (
	"context"
	"os"
	"os/signal"
	"sync"

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

	cartProvider, err := cart.NewCartRepository(ctx, wg, config.DbUser, config.DbPassword, config.DbHost, config.DbName)
	if err != nil {
		log.Fatal().Err(err).Msg("failed init cart pg repository")
		return
	}
	wg.Add(1)

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

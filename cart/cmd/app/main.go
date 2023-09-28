package main

import (
	"context"
	"os"
	"os/signal"
	"sync"

	grpcclients "route256/cart/internal/app/grpc/clients"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"route256/cart/internal/app"
	grpcserver "route256/cart/internal/app/grpc/server"
	httpclients "route256/cart/internal/app/http/clients"
	httpserver "route256/cart/internal/app/http/server"
	"route256/cart/internal/app/services"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	config := app.BuildConfig()

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	lomsService, err := grpcclients.NewLomsGrpcClient(ctx, wg, config.LomsGrpcServiceUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("failed init grpc loms client")
	}

	productService := httpclients.NewProductHttpClient(config.ProductServiceUrl)

	cartProvider := services.NewMockCartProvider(nil)
	cartService := services.NewCartService(productService, lomsService, cartProvider)

	httpServer := httpserver.NewHTTPServer(config.HttpServerAddress, cartService)
	grpcServer := grpcserver.NewGRPCServer(config.GrpcServerAddress, cartService)

	wg.Add(2)

	go httpServer.Run(ctx, wg)
	go grpcServer.Run(ctx, wg)

	wg.Wait()

	log.Info().Msg("goodbye")
}

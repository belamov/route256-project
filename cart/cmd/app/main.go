package main

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"route256/cart/internal/app/services"

	"route256/cart/internal/app"

	"route256/cart/internal/app/http/clients"
	"route256/cart/internal/app/http/server"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	config := app.BuildConfig()

	productService := clients.NewProductHttpClient(config.ProductServiceUrl)
	lomsService := clients.NewLomsHttpClient(config.LomsServiceUrl)
	cartProvider := services.NewMockCartProvider(nil)
	cartService := services.NewCartService(productService, lomsService, cartProvider)

	srv := server.NewHTTPServer(config.ServerAddress, cartService)

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	wg := &sync.WaitGroup{}

	wg.Add(1)

	srv.Run(ctx, wg)

	wg.Wait()

	log.Info().Msg("goodbye")
}

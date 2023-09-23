package main

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"route256/cart/internal/app/http/clients"
	"route256/cart/internal/app/http/server"

	"route256/cart/internal/app/domain/services"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"route256/cart/internal/app"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()
	// TODO: make configurable via env
	productService := clients.NewProductHttpClient("http://route256.pavl.uk:8080/get_product")
	lomsService := clients.NewLomsHttpClient("http://localhost:8083")
	// TODO: replace mock with real provider
	cartProvider := services.NewMockCartProvider(nil)
	cartService := services.NewCartService(productService, lomsService, cartProvider)

	config := app.BuildServerConfig()

	srv := server.NewHTTPServer(config.Address, cartService)

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	wg := &sync.WaitGroup{}

	wg.Add(2)

	srv.Run(ctx, wg)
	wg.Done()
	log.Info().Msg("goodbye")
}

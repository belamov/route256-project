package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"route256/cart/internal/app"
	mocks "route256/cart/internal/app/mocks"
	"route256/cart/internal/app/server"
	"route256/cart/internal/app/services"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()
	// TODO: make configurable via env
	productService := services.NewProductHttpClient("http://route256.pavl.uk:8080/get_product")
	lomsService := services.NewLomsHttpClient("http://localhost:8083")
	// TODO: replace mock with real provider
	cartProvider := mocks.NewMockCartProvider(nil)
	cartService := services.NewCartService(productService, lomsService, cartProvider)

	config := app.BuildServerConfig()

	srv := server.NewHTTPServer(config.Address, cartService)
	srv.Run()
}

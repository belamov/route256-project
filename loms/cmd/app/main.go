package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"route256/loms/internal/app"
	mocks "route256/loms/internal/app/mocks"
	"route256/loms/internal/app/server"
	"route256/loms/internal/app/services"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	orderProvider := mocks.NewMockOrdersProvider(nil)
	stocksProvider := mocks.NewMockStocksProvider(nil)
	lomsService := services.NewLomsService(orderProvider, stocksProvider)

	config := app.BuildServerConfig()

	srv := server.NewHTTPServer(config.Address, lomsService)
	srv.Run()
}

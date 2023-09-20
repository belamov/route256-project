package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	//cartService := services.NewCartService()
	//
	//config := app.BuildServerConfig()
	//
	//srv := server.NewHTTPServer(config.Address, cartService)
	//srv.Run()
}

package app

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	HttpServerAddress  string `default:"0.0.0.0:8080"`
	ProductServiceUrl  string `default:"http://route256.pavl.uk:8080/get_product"`
	LomsHttpServiceUrl string `default:"http://localhost:8080"`
	GrpcServerAddress  string `default:"localhost:8083"`
	LomsGrpcServiceUrl string `default:"localhost:8083"`
}

func BuildConfig() *Config {
	var config Config

	err := envconfig.Process("cart", &config)
	if err != nil {
		log.Panic().Err(err).Msg("cant build config")
	}

	log.Info().Msg("App config:\n" + config.String())

	return &config
}

func (config Config) String() string {
	return fmt.Sprintf(
		"HttpServerAddress: %v\n"+
			"ProductServiceUrl: %v\n"+
			"LomsHttpServiceUrl: %v\n"+
			"LomsGrpcServiceUrl: %v\n"+
			"GrpcServerAddress: %v\n",
		config.HttpServerAddress,
		config.ProductServiceUrl,
		config.LomsHttpServiceUrl,
		config.LomsGrpcServiceUrl,
		config.GrpcServerAddress,
	)
}

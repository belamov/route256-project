package app

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	HttpServerAddress        string `default:"0.0.0.0:8080"`
	ProductHttpServiceUrl    string `default:"http://route256.pavl.uk:8080/get_product"`
	ProductGrpcServiceUrl    string `default:"route256.pavl.uk:8082"`
	LomsHttpServiceUrl       string `default:"http://localhost:8080"`
	GrpcServerAddress        string `default:"localhost:8083"`
	GrpcGatewayServerAddress string `default:"0.0.0.0:8084"`
	LomsGrpcServiceUrl       string `default:"localhost:8083"`
	DbUser                   string `default:"postgres"`
	DbPassword               string `default:"password"`
	DbHost                   string `default:"db:5432"`
	DbName                   string `default:"cart"`
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
			"ProductHttpServiceUrl: %v\n"+
			"ProductGrpcServiceUrl: %v\n"+
			"LomsHttpServiceUrl: %v\n"+
			"LomsGrpcServiceUrl: %v\n"+
			"GrpcGatewayServerAddress: %v\n"+
			"GrpcServerAddress: %v\n"+
			"DbHost: %v\n"+
			"DbName: %v\n"+
			"DbUser: %v\n"+
			"DbPassword: %v\n",
		config.HttpServerAddress,
		config.ProductHttpServiceUrl,
		config.ProductGrpcServiceUrl,
		config.LomsHttpServiceUrl,
		config.LomsGrpcServiceUrl,
		config.GrpcGatewayServerAddress,
		config.GrpcServerAddress,
		config.DbHost,
		config.DbName,
		config.DbUser,
		config.DbPassword,
	)
}

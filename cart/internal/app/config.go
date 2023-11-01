package app

import (
	"github.com/rs/zerolog"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	HttpServerAddress         string        `default:"0.0.0.0:8080" split_words:"true"`
	ProductHttpServiceUrl     string        `default:"http://route256.pavl.uk:8080/get_product" split_words:"true"`
	TargetRpsToProductService int           `default:"10" split_words:"true"`
	ProductGrpcServiceUrl     string        `default:"route256.pavl.uk:8082" split_words:"true"`
	LomsHttpServiceUrl        string        `default:"http://localhost:8080" split_words:"true"`
	GrpcServerAddress         string        `default:"localhost:8083" split_words:"true"`
	GrpcGatewayServerAddress  string        `default:"0.0.0.0:8084" split_words:"true"`
	LomsGrpcServiceUrl        string        `default:"localhost:8083" split_words:"true"`
	DbUser                    string        `default:"postgres" split_words:"true"`
	DbPassword                string        `default:"password" split_words:"true"`
	DbHost                    string        `default:"db:5432" split_words:"true"`
	DbName                    string        `default:"cart" split_words:"true"`
	RedisAddress              string        `default:"redis:6379" split_words:"true"`
	LogLevel                  zerolog.Level `default:"3" split_words:"true"`
}

func BuildConfig() *Config {
	var config Config

	err := envconfig.Process("cart", &config)
	if err != nil {
		log.Panic().Err(err).Msg("cant build config")
	}

	log.Debug().Any("App config", config).Msg("config")

	return &config
}

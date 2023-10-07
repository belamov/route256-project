package app

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	HttpServerAddress          string        `default:"0.0.0.0:8080"`
	GrpcServerAddress          string        `default:"0.0.0.0:8083"`
	GrpcGatewayServerAddress   string        `default:"0.0.0.0:8084"`
	AllowedOrderUnpaidTime     time.Duration `default:"10m"`
	CancelUnpaidOrdersInterval time.Duration `default:"1m"`
	DbUser                     string        `default:"postgres" split_words:"true"`
	DbPassword                 string        `default:"password" split_words:"true"`
	DbHost                     string        `default:"db:5432" split_words:"true"`
	DbName                     string        `default:"loms" split_words:"true"`
}

func BuildConfig() *Config {
	var config Config

	err := envconfig.Process("loms", &config)
	if err != nil {
		log.Panic().Err(err).Msg("cant build config")
	}

	log.Info().Msg("App config:\n" + config.String())

	return &config
}

func (config Config) String() string {
	return fmt.Sprintf(
		"HttpServerAddress: %v\n"+
			"GrpcServerAddress: %v\n"+
			"GrpcGatewayServerAddress: %v\n"+
			"AllowedOrderUnpaidTime: %v\n"+
			"CancelUnpaidOrdersInterval: %v\n"+
			"DbHost: %v\n"+
			"DbName: %v\n"+
			"DbUser: %v\n"+
			"DbPassword: %v\n",
		config.HttpServerAddress,
		config.GrpcServerAddress,
		config.GrpcGatewayServerAddress,
		config.AllowedOrderUnpaidTime,
		config.CancelUnpaidOrdersInterval,
		config.DbHost,
		config.DbName,
		config.DbUser,
		config.DbPassword,
	)
}

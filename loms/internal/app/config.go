package app

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	HttpServerAddress          string        `default:"0.0.0.0:8080" split_words:"true"`
	GrpcServerAddress          string        `default:"0.0.0.0:8083" split_words:"true"`
	GrpcGatewayServerAddress   string        `default:"0.0.0.0:8084" split_words:"true"`
	AllowedOrderUnpaidTime     time.Duration `default:"10m" split_words:"true"`
	CancelUnpaidOrdersInterval time.Duration `default:"1m" split_words:"true"`
	DbUser                     string        `default:"postgres" split_words:"true"`
	DbPassword                 string        `default:"password" split_words:"true"`
	DbHost                     string        `default:"localhost:5432" split_words:"true"`
	DbName                     string        `default:"loms" split_words:"true"`
	KafkaBrokers               []string      `default:"localhost:9091,localhost:9092,localhost:9093" split_words:"true"`
	OutboxId                   string        `default:"notifications-1" split_words:"true"`
	OutboxSendInterval         time.Duration `default:"1m" split_words:"true"`
	OutboxRetryInterval        time.Duration `default:"10m" split_words:"true"`
}

func BuildConfig() *Config {
	var config Config

	err := envconfig.Process("loms", &config)
	if err != nil {
		log.Panic().Err(err).Msg("cant build config")
	}

	log.Info().Any("App config", config)

	return &config
}

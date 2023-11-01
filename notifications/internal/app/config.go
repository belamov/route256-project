package app

import (
	"github.com/rs/zerolog"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	KafkaBrokers    []string      `default:"localhost:9091,localhost:9092,localhost:9093" split_words:"true"`
	TopicNames      []string      `default:"order-status-changed" split_words:"true"`
	ConsumerGroupId string        `default:"notifications" split_words:"true"`
	LogLevel        zerolog.Level `default:"3" split_words:"true"`
}

func BuildConfig() *Config {
	var config Config

	err := envconfig.Process("loms", &config)
	if err != nil {
		log.Panic().Err(err).Msg("cant build config")
	}

	log.Debug().Any("App config", config).Msg("config")

	return &config
}

package app

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	KafkaBrokers    []string `default:"localhost:9091,localhost:9092,localhost:9093" split_words:"true"`
	TopicNames      []string `default:"order-status-changed" split_words:"true"`
	ConsumerGroupId string   `default:"notifications" split_words:"true"`
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
		"KafkaBrokers: %v\n"+
			"TopicNames: %v\n"+
			"ConsumerGroupId: %v\n",
		config.KafkaBrokers,
		config.TopicNames,
		config.ConsumerGroupId,
	)
}

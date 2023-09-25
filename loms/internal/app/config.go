package app

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	ServerAddress              string        `default:"0.0.0.0:8080"`
	AllowedOrderUnpaidTime     time.Duration `default:"1m"`
	CancelUnpaidOrdersInterval time.Duration `default:"10m"`
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
		"ServerAddress: %v\n"+
			"AllowedOrderUnpaidTime: %v\n"+
			"CancelUnpaidOrdersInterval: %v\n",
		config.ServerAddress,
		config.AllowedOrderUnpaidTime,
		config.CancelUnpaidOrdersInterval,
	)
}

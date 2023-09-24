package app

import (
	"flag"
	"os"
	"time"
)

type ServerConfig struct {
	Address                    string
	AllowedOrderUnpaidTime     time.Duration
	CancelUnpaidOrdersInterval time.Duration
}

func BuildServerConfig() *ServerConfig {
	defaultAddress := "0.0.0.0:8080"
	flag.Parse()

	cfg := &ServerConfig{
		Address:                    coalesceStrings(os.Getenv("ADDRESS"), defaultAddress),
		CancelUnpaidOrdersInterval: time.Second,
	}

	return cfg
}

func coalesceStrings(strings ...string) string {
	for _, str := range strings {
		if str != "" {
			return str
		}
	}
	return ""
}

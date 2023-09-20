package app

import (
	"flag"
	"os"
)

type ServerConfig struct {
	Address string
}

func BuildServerConfig() *ServerConfig {
	defaultAddress := "0.0.0.0:8080"
	flag.Parse()

	cfg := &ServerConfig{
		Address: coalesceStrings(os.Getenv("ADDRESS"), defaultAddress),
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

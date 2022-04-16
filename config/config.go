package config

import (
	"log"
)

type Config struct {
	Addr   string
	Logger *log.Logger
}

func NewConfig(logger *log.Logger) *Config {
	return &Config{
		Addr:   ":8080",
		Logger: logger,
	}
}

package config

import (
	"log"
	"os"
)

type Config struct {
	Addr   string
	Logger *log.Logger
}

func NewConfig() *Config {
	return &Config{
		Addr:   ":8080",
		Logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	Addr                 string
	DBSnapshotPath       string
	ShortURLDomain       string
	Logger               *log.Logger
	SnapshotSaveInterval time.Duration
}

const defaultAddr = ":8080"
const defaultShortURLDomain = "http://localhost:8080"

func NewConfig(logger *log.Logger) *Config {
	addr := os.Getenv("APP_ADDR")
	shortURLDomain := os.Getenv("SHORT_URL_DOMAIN")

	if addr == "" {
		addr = defaultAddr
	}

	if shortURLDomain == "" {
		shortURLDomain = defaultShortURLDomain
	}

	return &Config{
		Addr:                 addr,
		ShortURLDomain:       shortURLDomain,
		Logger:               logger,
		DBSnapshotPath:       "snapshot.db",
		SnapshotSaveInterval: 5 * time.Second,
	}
}

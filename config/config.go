package config

import (
	"log"
	"time"
)

type Config struct {
	Addr                 string
	DBSnapshotPath       string
	ShortURLDomain       string
	Logger               *log.Logger
	SnapshotSaveInterval time.Duration
}

func NewConfig(logger *log.Logger) *Config {
	return &Config{
		Addr:                 ":8080",
		ShortURLDomain:       "http://localhost:8080",
		Logger:               logger,
		DBSnapshotPath:       "./internal/platform/snapshot/snapshot.db",
		SnapshotSaveInterval: 5 * time.Second,
	}
}

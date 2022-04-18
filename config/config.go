package config

import (
	"log"
	"time"
)

type Config struct {
	Addr                 string
	Logger               *log.Logger
	DBSnapshotPath       string
	SnapshotSaveInterval time.Duration
}

func NewConfig(logger *log.Logger) *Config {
	return &Config{
		Addr:                 ":8080",
		Logger:               logger,
		DBSnapshotPath:       "./db/snapshot.db",
		SnapshotSaveInterval: 5 * time.Second,
	}
}

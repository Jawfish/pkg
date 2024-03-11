package main

import (
	"log/slog"
)

type Config struct {
	packageManager string
	architecture   string
}

func NewConfig(packageManager string, architecture string) *Config {
	slog.Info("creating new config", "packageManager", packageManager, "architecture", architecture)

	return &Config{packageManager: packageManager, architecture: architecture}

}

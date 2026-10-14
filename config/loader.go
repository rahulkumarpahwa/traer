package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

func Load(path string) (*Config, error) {
	var cfg Config

	// If YAML exists → load it
	if path != "" {
		if _, err := os.Stat(path); err == nil {
			if err := cleanenv.ReadConfig(path, cfg); err != nil {
				return nil, fmt.Errorf("failed to read yaml config: %w", err)
			}
		}
	}

	return &cfg, nil
}

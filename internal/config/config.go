package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	ENV string `env:"ENV" envDefault:"development"`
	API struct {
		PORT int `env:"PORT" envDefault:"8080"`
	}
	DB struct {
		PGPORT     int    `env:"PGPORT" envDefault:"5432"`
		PGDATABASE string `env:"PGDATABASE" envDefault:"local"`
		PGUSER     string `env:"PGUSER" envDefault:"root"`
		PGPASSWORD string `env:"PGPASSWORD" envDefault:""`
	}
}

func LoadConfig() (*Config, error) {
	var cfg Config
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

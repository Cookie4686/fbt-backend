package util

import (
	"fmt"
	"os"
	"testing"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	ENV string `env:"ENV"`
	API struct {
		PORT int `env:"PORT"`
	}
	DB struct {
		PGPORT     int    `env:"PGPORT"`
		PGDATABASE string `env:"PGDATABASE"`
		PGUSER     string `env:"PGUSER"`
		PGPASSWORD string `env:"PGPASSWORD"`
	}
	ENCRYPTION_KEY string `env:"ENCRYPTION_KEY"`
	MAIL           struct {
		ADDRESS  string `env:"MAIL_ADDRESS"`
		SERVER   string `env:"MAIL_SERVER"`
		PORT     int    `env:"MAIL_PORT"`
		USERNAME string `env:"MAIL_USERNAME"`
		PASSWORD string `env:"MAIL_PASSWORD"`
	}
}

func NewConfig() (*Config, error) {
	environment, ok := os.LookupEnv("ENV")
	if !ok {
		if testing.Testing() {
			environment = "test"
		} else {
			environment = "production"
		}
	}

	if err := godotenv.Load(fmt.Sprintf(".env.%s", environment)); err != nil {
		if err := godotenv.Load(".env"); err != nil {
			godotenv.Load(".env.example")
		}
	}

	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	} else {
		return &cfg, nil
	}
}

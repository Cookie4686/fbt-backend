package util

import (
	"errors"
	"io/fs"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	ENV string `env:"ENV" envDefault:"development"`
	API struct {
		PORT int `env:"PORT" envDefault:"8080"`
	}
	DB struct {
		PGPORT     int    `env:"PGPORT"     envDefault:"5432"`
		PGDATABASE string `env:"PGDATABASE" envDefault:"local"`
		PGUSER     string `env:"PGUSER"     envDefault:"root"`
		PGPASSWORD string `env:"PGPASSWORD" envDefault:""`
	}
	ENCRYPTION_KEY string `env:"ENCRYPTION_KEY" envDefault:""`
	MAIL           struct {
		ADDRESS  string `env:"MAIL_ADDRESS"  envDefault:""`
		SERVER   string `env:"MAIL_SERVER"   envDefault:""`
		PORT     int    `env:"MAIL_PORT"     envDefault:"25"`
		USERNAME string `env:"MAIL_USERNAME" envDefault:""`
		PASSWORD string `env:"MAIL_PASSWORD" envDefault:""`
	}
}

func NewConfig() (*Config, error) {
	environment, ok := os.LookupEnv("ENV")
	var file string
	if ok {
		switch environment {
		case "development":
			file = ".env.development"
		case "test":
			file = ".env.test"
		case "production":
			file = ".env.production"
		}

		if err := godotenv.Load(file); err != nil && errors.Is(err, &fs.PathError{}) {
			return nil, err
		}
	}
	if err := godotenv.Load(".env"); err != nil && errors.Is(err, &fs.PathError{}) {
		return nil, err
	}

	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	} else {
		return &cfg, nil
	}
}

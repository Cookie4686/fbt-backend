package util

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wneessen/go-mail"
	"go.uber.org/zap"
)

type Dependency struct {
	Logger *zap.Logger
	CFG    *Config
	DB     *pgxpool.Pool
	Mail   *mail.Client
}

func NewDependency(ctx context.Context, envFilePath string) (*Dependency, error) {
	if cfg, err := NewConfig(envFilePath); err != nil {
		return nil, err
	} else if logger, err := NewLogger(cfg); err != nil {
		return nil, err
	} else if db, err := NewDatabasePool(ctx, cfg); err != nil {
		return nil, err
	} else if mail, err := NewMailClient(cfg); err != nil {
		return nil, err
	} else {
		return &Dependency{Logger: logger, CFG: cfg, DB: db, Mail: mail}, nil
	}
}

func NewLogger(cfg *Config) (logger *zap.Logger, err error) {
	if cfg.ENV == "" || cfg.ENV == "development" {
		return zap.NewDevelopment()
	} else {
		return zap.NewProduction()
	}
}

func NewDatabasePool(ctx context.Context, cfg *Config) (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(fmt.Sprintf(
		"user=%v password=%v port=%v dbname=%v",
		cfg.DB.PGUSER, cfg.DB.PGPASSWORD, cfg.DB.PGPORT, cfg.DB.PGDATABASE),
	)
	if err != nil {
		return nil, err
	}

	conn, err := pgxpool.NewWithConfig(ctx, dbConfig)
	return conn, err
}

func NewMailClient(cfg *Config) (*mail.Client, error) {
	client, err := mail.NewClient(
		cfg.MAIL.SERVER,
		mail.WithPort(cfg.MAIL.PORT),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(cfg.MAIL.USERNAME),
		mail.WithPassword(cfg.MAIL.PASSWORD),
	)
	if err != nil {
		return nil, err
	}
	return client, err
}

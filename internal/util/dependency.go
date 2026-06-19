package util

import (
	"context"
	"fbt/backend/internal/config"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Dependency struct {
	Logger *zap.Logger
	DB     *pgxpool.Pool
	CFG    *config.Config
}

func NewDependency(logger *zap.Logger, db *pgxpool.Pool, cfg *config.Config) *Dependency {
	return &Dependency{
		Logger: logger,
		DB:     db,
		CFG:    cfg,
	}
}

func NewLogger(cfg *config.Config) (logger *zap.Logger, err error) {
	if cfg.ENV == "" || cfg.ENV == "development" {
		return zap.NewDevelopment()
	} else {
		return zap.NewProduction(zap.Fields(
			zap.String("env", cfg.ENV),
		))
	}
}

func NewDatabasePool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
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

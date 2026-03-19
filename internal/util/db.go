package util

import (
	"context"
	"fbt/backend/internal/config"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

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

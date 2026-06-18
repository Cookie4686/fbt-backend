package dependency

import (
	"fbt/backend/internal/config"

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

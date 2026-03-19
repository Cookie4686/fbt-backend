package util

import (
	"fbt/backend/internal/config"

	"go.uber.org/zap"
)

func NewLogger(cfg *config.Config) (logger *zap.Logger, err error) {
	if cfg.ENV == "" || cfg.ENV == "development" {
		return zap.NewDevelopment()
	} else {
		return zap.NewProduction(zap.Fields(
			zap.String("env", cfg.ENV),
		))
	}
}

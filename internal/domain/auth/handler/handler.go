package handler

import (
	"fbt/backend/internal/config"
	"fbt/backend/internal/domain/auth/repo"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type AuthHandler struct {
	logger *zap.Logger
	cfg    *config.Config
	repo   *repo.AuthRepository
}

func NewAuthHandler(logger *zap.Logger, db *pgxpool.Pool, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		logger: logger,
		cfg:    cfg,
		repo:   repo.NewAuthRepository(db, cfg),
	}
}

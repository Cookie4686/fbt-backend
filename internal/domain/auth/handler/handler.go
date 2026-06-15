package handler

import (
	"fbt/backend/internal/domain/auth/repo"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type AuthHandler struct {
	logger *zap.Logger
	repo   *repo.AuthRepository
}

func NewAuthHandler(logger *zap.Logger, db *pgxpool.Pool) *AuthHandler {
	return &AuthHandler{
		logger: logger,
		repo:   repo.NewAuthRepository(db),
	}
}

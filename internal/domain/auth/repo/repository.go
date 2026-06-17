package repo

import (
	"fbt/backend/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	db  *pgxpool.Pool
	cfg *config.Config
}

func NewAuthRepository(db *pgxpool.Pool, cfg *config.Config) *AuthRepository {
	return &AuthRepository{db, cfg}
}

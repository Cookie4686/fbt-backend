package service

import (
	"context"
	"fbt/backend/internal/dependency"
	"fbt/backend/internal/domain/auth/model"
)

type Service interface {
	CreateSession(ctx context.Context, userId string) (*model.Session, error)
	Validate(ctx context.Context, sessionId string) (*model.Auth, error)
	UpdateSessionExpiration(context.Context, *model.Session) error
	InvalidateSession(context.Context, *model.Session) error

	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)

	Decrypt(encryptedValue string) (*string, error)
	Encrypt(value string) (*string, error)
}

type service struct {
	*dependency.Dependency
}

func NewService(d *dependency.Dependency) Service {
	return Service(&service{Dependency: d})
}

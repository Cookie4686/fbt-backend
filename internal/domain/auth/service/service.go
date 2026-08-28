// Package service for common-use services in authentication domain
package service

import (
	"context"

	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/util"
)

type Service interface {
	CreateSession(ctx context.Context, userID string, twoFactorVerified bool) (*model.Session, error)
	Validate(ctx context.Context, sessionID string) (*model.Auth, error)
	UpdateSessionExpiration(context.Context, *model.Session) error
	InvalidateSession(context.Context, *model.Session) error

	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)

	SendVerificationMail(email, otp string) error

	GenerateSalt() []byte
	HashPassword(password string, salt []byte) []byte

	Decrypt(encryptedValue string) (*string, error)
	Encrypt(value string) (*string, error)
}

type service struct {
	*util.Dependency
}

func NewService(d *util.Dependency) Service {
	return Service(&service{Dependency: d})
}

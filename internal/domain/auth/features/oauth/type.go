package oauth

import (
	"context"
	"fbt/backend/internal/domain/auth/model"
	"net/http"
)

type Controller interface {
	Register(context.Context, *RegisterPayload) (*RegisterResponse, error)
	Login(context.Context, *LoginPayload) (*LoginResponse, error)
	Status(context.Context, *model.Auth) (*StatusResponse, error)
}

type Feature struct {
	Register    http.Handler
	Login       http.Handler
	AUTH_Status http.Handler
}

type Repo interface {
	GetUserOAuth(ctx context.Context, provider string, idToken string) (*model.UserOAuth, error)
	LinkOAuth(ctx context.Context, provider string, userID string, idToken string) error
	CreateOAuthRegistration(ctx context.Context, provider string, oauthRegistration *model.OauthRegistration) error
	GetOAuthRegistration(ctx context.Context, registrationId string) (*model.OauthRegistration, error)
	DeleteOAuthRegistration(ctx context.Context, provider string, idToken string) error
	OAuthRegister(ctx context.Context, registrationId string, user *model.User, session *model.Session) error
	GetUserProvider(ctx context.Context, userID string) ([]string, error)
}

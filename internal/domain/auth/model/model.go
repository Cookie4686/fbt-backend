package model

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const SessionExpiresIn = 24 * time.Hour

type User struct {
	Id              string      `json:"id" db:"user_id"`
	Username        string      `json:"username" db:"username"`
	Email           string      `json:"email" db:"email"`
	EmailVerified   bool        `json:"emailVerified" db:"email_verified"`
	Password        pgtype.Text `json:"-" db:"password"`
	PasswordSalt    pgtype.Text `json:"-" db:"password_salt"`
	PasswordEnabled bool        `json:"password_enabled" db:"password_enabled"`
}

type Session struct {
	Id                string    `json:"id" db:"session_id"`
	UserId            string    `json:"userId" db:"user_id"`
	ExpiresAt         time.Time `json:"expiresAt" db:"expires_at"`
	TwoFactorVerified bool      `json:"twoFactorVerified" db:"two_factor_verified"`
}

type Auth struct {
	Session Session `json:"session"`
	User    User    `json:"user"`
}

type OauthProvider struct {
	OauthProviderID int16       `db:"oauth_provider_id"`
	Name            pgtype.Text `db:"name"`
}

type UserOAuth struct {
	UserID     string `json:"user_id" db:"user_id"`
	ProviderID int16  `json:"provider_id" db:"oauth_provider_id"`
	IDToken    string `json:"id_token" db:"id_token"`
}

type OauthRegistration struct {
	RegistrationID  string    `json:"registration_id" db:"registration_id"`
	IDToken         string    `json:"id_token" db:"id_token"`
	OauthProviderID int16     `json:"oauth_provider_id" db:"oauth_provider_id"`
	ExpiresAt       time.Time `json:"expiresAt" db:"expires_at"`
}

// type UserRepository interface {
// 	Get(ctx context.Context, userId string) (*User, error)
// 	Create(ctx context.Context, _ *User) (string, error)
// }

// type SessionRepository interface {
// 	Validate(ctx context.Context, sessionId string) (SessionContext, error)
// 	Create(ctx context.Context, userId string) (string, error)
// }

package model

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type OauthProvider struct {
	OauthProviderID int16       `db:"oauth_provider_id"`
	Name            pgtype.Text `db:"name"`
}

type UserOAuth struct {
	UserID     string `db:"user_id"`
	ProviderID int16  `db:"oauth_provider_id"`
	IDToken    string `db:"id_token"`
}

type OauthRegistration struct {
	RegistrationID  string    `db:"registration_id"`
	IDToken         string    `db:"id_token"`
	OauthProviderID int16     `db:"oauth_provider_id"`
	EmailVerified   bool      `db:"email_verified"`
	ExpiresAt       time.Time `db:"expires_at"`
}

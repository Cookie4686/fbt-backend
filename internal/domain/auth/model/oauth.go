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

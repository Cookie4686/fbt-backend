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

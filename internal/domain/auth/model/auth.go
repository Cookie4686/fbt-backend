package model

import (
	authv1 "fbt/backend/gen/proto/go/auth/v1"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const SessionExpiresIn = 24 * time.Hour
const SessionMaxAge = 72 * time.Hour

type User struct {
	Id              string      `db:"user_id"`
	Username        string      `db:"username"`
	Email           string      `db:"email"`
	EmailVerified   bool        `db:"email_verified"`
	Password        pgtype.Text `db:"password"`
	PasswordSalt    pgtype.Text `db:"password_salt"`
	PasswordEnabled bool        `db:"password_enabled"`
}

func (u *User) ToProto() *authv1.User {
	return &authv1.User{
		Id:              u.Id,
		Username:        u.Username,
		Email:           u.Email,
		EmailVerified:   u.EmailVerified,
		PasswordEnabled: u.PasswordEnabled,
	}
}

type Session struct {
	Id                string    `db:"session_id"`
	UserId            string    `db:"user_id"`
	CreatedAt         time.Time `db:"created_at"`
	ExpiresAt         time.Time `db:"expires_at"`
	TwoFactorVerified bool      `db:"two_factor_verified"`
}

func (s *Session) ToProto() *authv1.Session {
	return &authv1.Session{
		Id:                s.Id,
		UserId:            s.UserId,
		TwoFactorVerified: s.TwoFactorVerified,
		ExpiresAt:         timestamppb.New(s.ExpiresAt),
	}
}

type Auth struct {
	Session Session
	User    User
}

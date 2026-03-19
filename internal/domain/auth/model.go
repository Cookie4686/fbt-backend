package auth

import "time"

type User struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Password      string `json:"password"`
	PasswordSalt  string `json:"password_salt"`
}

type Session struct {
	Id                string    `json:"id"`
	UserId            string    `json:"userId"`
	ExpiresAt         time.Time `json:"expiresAt"`
	TwoFactorVerified bool      `json:"twoFactorVerified"`
}

// type UserRepository interface {
// 	Get(ctx context.Context, userId string) (*User, error)
// 	Create(ctx context.Context, _ *User) (string, error)
// }

// type SessionRepository interface {
// 	Validate(ctx context.Context, sessionId string) (SessionContext, error)
// 	Create(ctx context.Context, userId string) (string, error)
// }

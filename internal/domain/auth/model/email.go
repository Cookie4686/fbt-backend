package model

import "github.com/jackc/pgx/v5/pgtype"

type EmailVerification struct {
	UserID         string           `db:"user_id"`
	VerificationID string           `db:"verification_id"`
	Otp            string           `db:"otp"`
	ExpiresAt      pgtype.Timestamp `db:"expires_at"`
}

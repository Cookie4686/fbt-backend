package email

import (
	"context"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/util"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo interface {
	GetEmailVerification(ctx context.Context, userID string) (*model.EmailVerification, error)
	CreateEmailVerification(context.Context, *model.EmailVerification) error
	VerifyEmail(ctx context.Context, userID string) error
	DeleteEmailVerification(ctx context.Context, userID string) error
}

type repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) Repo {
	return Repo(&repo{db})
}

func (s *repo) GetEmailVerification(ctx context.Context, userID string) (*model.EmailVerification, error) {
	query := `
		SELECT * FROM email_verification
		WHERE email_verification.user_id = @user_id
	`
	args := pgx.NamedArgs{"user_id": userID}

	emailVerification, err := util.FetchOne[model.EmailVerification](s.db, ctx, query, args)
	if err != nil {
		return nil, err
	}

	return emailVerification, nil
}

func (s *repo) CreateEmailVerification(ctx context.Context, email *model.EmailVerification) error {
	query := `
		INSERT INTO email_verification(user_id, verification_id, otp, expires_at)
		VALUES (@user_id, @verification_id, @otp, @expires_at)
	`
	args := pgx.NamedArgs{
		"user_id":         email.UserID,
		"verification_id": email.VerificationID,
		"otp":             email.Otp,
		"expires_at":      email.ExpiresAt,
	}
	_, err := s.db.Exec(ctx, query, args)

	return err
}

func (s *repo) VerifyEmail(ctx context.Context, userID string) error {
	batch := &pgx.Batch{}
	batch.Queue(`
		DELETE FROM email_verification
		WHERE user_id = @user_id
	`,
		pgx.NamedArgs{"user_id": userID},
	)
	batch.Queue(`
		UPDATE users
		SET email_verified = true
		WHERE user_id = @user_id
	`, pgx.NamedArgs{"user_id": userID},
	)

	_, err := s.db.SendBatch(ctx, batch).Exec()

	return err
}

func (s *repo) DeleteEmailVerification(ctx context.Context, userID string) error {
	query := `
		DELETE FROM email_verification
		WHERE user_id = @user_id
	`
	args := pgx.NamedArgs{"user_id": userID}
	_, err := s.db.Exec(ctx, query, args)

	return err
}

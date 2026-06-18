package mfa

import (
	"context"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	db *pgxpool.Pool
}

func newRepo(db *pgxpool.Pool) *repo {
	return &repo{db}
}

func (s *repo) GetMFAList(ctx context.Context, userID string) (*model.MfaList, error) {
	query := `
		SELECT mfa_totp.id FROM users
		LEFT JOIN mfa_totp ON users.user_id = mfa_totp.user_id
		WHERE users.user_id = @user_id
	`
	args := pgx.NamedArgs{"user_id": userID}
	row := s.db.QueryRow(ctx, query, args)

	var list struct {
		TotpID *string
	}
	err := row.Scan(&list.TotpID)
	if err != nil {
		return nil, err
	}

	return &model.MfaList{
		Totp: list.TotpID != nil,
	}, nil
}

func (s *repo) GetTOTP(ctx context.Context, userID string) (*model.MfaTotp, error) {
	query := `
		SELECT * FROM mfa_totp
		WHERE user_id = @user_id
	`
	args := pgx.NamedArgs{"user_id": userID}
	totp, err := s.fetchTOTP(ctx, query, args)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound
		}
		return nil, err
	}
	return totp, nil
}

func (s *repo) UpsertTOTP(ctx context.Context, key string, userID string) error {
	query := `
		INSERT INTO mfa_totp(key, user_id)
		VALUES (@key, @user_id)
		ON CONFLICT(user_id)
		DO UPDATE SET
			key = EXCLUDED.key
	`
	args := pgx.NamedArgs{"key": key, "user_id": userID}
	_, err := s.db.Exec(ctx, query, args)
	return err
}

func (s *repo) fetchTOTP(ctx context.Context, query string, args ...any) (*model.MfaTotp, error) {
	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	totp, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.MfaTotp])
	if err != nil {
		return nil, err
	}
	return &totp, nil
}

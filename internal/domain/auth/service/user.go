package service

import (
	"context"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/errors"

	"github.com/jackc/pgx/v5"
)

func (s *AuthService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `SELECT * FROM users WHERE username = @username`
	args := pgx.NamedArgs{"username": username}
	user, err := s.fetchUser(ctx, query, args)
	if err != nil {
		if err == pgx.ErrNoRows {
			return user, errors.NotFound
		}
		return nil, err
	}
	return user, err
}

func (s *AuthService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT * FROM users WHERE email = @email`
	args := pgx.NamedArgs{"email": email}
	user, err := s.fetchUser(ctx, query, args)
	if err != nil {
		if err == pgx.ErrNoRows {
			return user, errors.NotFound
		}
		return nil, err
	}
	return user, err
}

func (s *AuthService) fetchUser(ctx context.Context, query string, args ...any) (*model.User, error) {
	rows, err := s.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.User])
	if err != nil {
		return nil, err
	}
	return &user, nil
}

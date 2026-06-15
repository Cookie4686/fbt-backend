package repo

import (
	"context"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/errors"

	"github.com/jackc/pgx/v5"
)

func (s *AuthRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
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

func (s *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
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

func (s *AuthRepository) fetchUsers(ctx context.Context, query string, args ...any) ([]model.User, error) {
	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.User])
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *AuthRepository) fetchUser(ctx context.Context, query string, args ...any) (*model.User, error) {
	rows, err := s.db.Query(ctx, query, args...)
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

package service

import (
	"context"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/errors"
	"fbt/backend/internal/util"

	"github.com/jackc/pgx/v5"
)

func (s *service) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `SELECT * FROM users WHERE username = @username`
	args := pgx.NamedArgs{"username": username}
	user, err := util.FetchOne[model.User](s.DB, ctx, query, args)
	if err != nil {
		if err == pgx.ErrNoRows {
			return user, errors.NotFound
		}
		return nil, err
	}
	return user, err
}

func (s *service) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT * FROM users WHERE email = @email`
	args := pgx.NamedArgs{"email": email}
	user, err := util.FetchOne[model.User](s.DB, ctx, query, args)
	if err != nil {
		if err == pgx.ErrNoRows {
			return user, errors.NotFound
		}
		return nil, err
	}
	return user, err
}

package auth

import (
	"context"
	"fbt/backend/internal/errors"

	"github.com/jackc/pgx/v5"
)

func (s *AuthRepository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
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

func (s *AuthRepository) fetchUsers(ctx context.Context, query string, args ...any) ([]User, error) {
	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.Id,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.PasswordSalt,
			&user.EmailVerified,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *AuthRepository) fetchUser(ctx context.Context, query string, args ...any) (*User, error) {
	var user User
	err := s.db.QueryRow(ctx, query, args...).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.PasswordSalt,
		&user.EmailVerified,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

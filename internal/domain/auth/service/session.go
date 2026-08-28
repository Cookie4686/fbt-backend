package service

import (
	"context"
	"time"

	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/errors"

	"github.com/jackc/pgx/v5"
)

const sessionRenewInterval = model.SessionExpiresIn / 2

// SECURITY: Do we need to hash session id?

func (s *service) CreateSession(ctx context.Context, userID string, twoFactorVerified bool) (*model.Session, error) {
	session := model.NewSession(userID, twoFactorVerified)

	query := `
		INSERT INTO sessions(session_id, user_id, created_at, expires_at, two_factor_verified)
		VALUES (@sessionId, @userId, @createdAt, @expiresAt, @twoFactorVerified)
	`
	args := pgx.NamedArgs{
		"sessionId":         session.ID,
		"userId":            session.UserID,
		"createdAt":         session.CreatedAt,
		"expiresAt":         session.ExpiresAt,
		"twoFactorVerified": session.TwoFactorVerified,
	}

	if _, err := s.DB.Exec(ctx, query, args); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *service) Validate(ctx context.Context, sessionID string) (*model.Auth, error) {
	query := `
		SELECT session_id, sessions.user_id, created_at, expires_at, two_factor_verified, users.user_id, username, email, email_verified, password, password_salt, password_enabled
		FROM sessions
		LEFT JOIN users ON sessions.user_id = users.user_id
		WHERE session_id = @sessionId
	`
	args := pgx.NamedArgs{"sessionId": sessionID}
	row := s.DB.QueryRow(ctx, query, args)

	var auth model.Auth

	err := row.Scan(
		&auth.Session.ID,
		&auth.Session.UserID,
		&auth.Session.CreatedAt,
		&auth.Session.ExpiresAt,
		&auth.Session.TwoFactorVerified,
		&auth.User.ID,
		&auth.User.Username,
		&auth.User.Email,
		&auth.User.EmailVerified,
		&auth.User.Password,
		&auth.User.PasswordSalt,
		&auth.User.PasswordEnabled,
	)
	if err == pgx.ErrNoRows {
		return nil, errors.NotFound
	}

	if err != nil {
		return nil, err
	}

	if time.Now().After(auth.Session.ExpiresAt) {
		if err := s.InvalidateSession(ctx, &auth.Session); err != nil {
			return nil, errors.DBError
		}

		return nil, errors.SessionExpire
	}

	if time.Now().After(auth.Session.ExpiresAt.Add(-sessionRenewInterval)) {
		newExpiresAt := time.Now().Add(model.SessionExpiresIn)
		maxExpiresAt := auth.Session.CreatedAt.Add(model.SessionMaxAge)

		if newExpiresAt.After(maxExpiresAt) {
			auth.Session.ExpiresAt = maxExpiresAt
		} else {
			auth.Session.ExpiresAt = newExpiresAt
		}

		if err := s.UpdateSessionExpiration(ctx, &auth.Session); err != nil {
			return nil, errors.DBError
		}
	}

	return &auth, nil
}

func (s *service) UpdateSessionExpiration(ctx context.Context, session *model.Session) error {
	query := `
		UPDATE sessions
		SET expires_at = @expiresAt
		WHERE session_id = @sessionId
	`
	args := pgx.NamedArgs{"expiresAt": session.ExpiresAt, "sessionId": session.ID}
	_, err := s.DB.Exec(ctx, query, args)

	return err
}

func (s *service) InvalidateSession(ctx context.Context, session *model.Session) error {
	query := `
		DELETE FROM sessions
		WHERE session_id = @sessionId
	`
	args := pgx.NamedArgs{"sessionId": session.ID}
	_, err := s.DB.Exec(ctx, query, args)

	return err
}

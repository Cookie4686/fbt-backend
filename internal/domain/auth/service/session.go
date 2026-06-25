package service

import (
	"context"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/errors"
	"fbt/backend/internal/util"
	"time"

	"github.com/jackc/pgx/v5"
)

// SECURITY: Do we need to hash session id?

func (s *service) CreateSession(ctx context.Context, userId string, twoFactorVerified bool) (*model.Session, error) {
	session := &model.Session{
		Id:                util.GenerateBase64UUID(),
		UserId:            userId,
		ExpiresAt:         time.Now().Add(model.SessionExpiresIn),
		TwoFactorVerified: twoFactorVerified,
	}

	query := `
		INSERT INTO sessions(session_id, user_id, expires_at, two_factor_verified)
		VALUES (@sessionId, @userId, @expiresAt, @twoFactorVerified)
	`
	args := pgx.NamedArgs{
		"sessionId":         session.Id,
		"userId":            session.UserId,
		"expiresAt":         session.ExpiresAt,
		"twoFactorVerified": session.TwoFactorVerified,
	}

	if _, err := s.DB.Exec(ctx, query, args); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *service) Validate(ctx context.Context, sessionId string) (*model.Auth, error) {
	query := `
		SELECT session_id, sessions.user_id, expires_at, two_factor_verified, users.user_id, username, email, email_verified, password, password_salt
		FROM sessions
		LEFT JOIN users ON sessions.user_id = users.user_id
		WHERE session_id = @sessionId
	`
	args := pgx.NamedArgs{"sessionId": sessionId}
	row := s.DB.QueryRow(ctx, query, args)

	var auth model.Auth
	err := row.Scan(
		&auth.Session.Id,
		&auth.Session.UserId,
		&auth.Session.ExpiresAt,
		&auth.Session.TwoFactorVerified,
		&auth.User.Id,
		&auth.User.Username,
		&auth.User.Email,
		&auth.User.EmailVerified,
		&auth.User.Password,
		&auth.User.PasswordSalt,
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
	if time.Now().After(auth.Session.ExpiresAt.Add(-model.SessionExpiresIn / 2)) {
		auth.Session.ExpiresAt = time.Now().Add(model.SessionExpiresIn)
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
	args := pgx.NamedArgs{"expiresAt": session.ExpiresAt, "sessionId": session.Id}
	_, err := s.DB.Exec(ctx, query, args)
	return err
}

func (s *service) InvalidateSession(ctx context.Context, session *model.Session) error {
	query := `
		DELETE FROM sessions
		WHERE session_id = @sessionId
	`
	args := pgx.NamedArgs{"sessionId": session.Id}
	_, err := s.DB.Exec(ctx, query, args)
	return err
}

func (s *service) fetchSessions(ctx context.Context, query string, args ...any) ([]model.Session, error) {
	rows, err := s.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Session])
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (s *service) fetchSession(ctx context.Context, query string, args ...any) (*model.Session, error) {
	rows, err := s.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	session, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Session])
	if err != nil {
		return nil, err
	}
	return &session, nil
}

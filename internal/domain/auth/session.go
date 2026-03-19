package auth

import (
	"context"
	"fbt/backend/internal/util"
	"time"

	"github.com/jackc/pgx/v5"
)

const sessionExpiresIn = 30 * 24 * time.Hour

// SECURITY: Do we need to hash session id?

func (s *AuthRepository) CreateSession(ctx context.Context, userId string) (*Session, error) {
	session := &Session{
		Id:                util.GenerateBase64UUID(),
		UserId:            userId,
		ExpiresAt:         time.Now().Add(sessionExpiresIn),
		TwoFactorVerified: false,
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

	if _, err := s.db.Exec(ctx, query, args); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *AuthRepository) fetchSessions(ctx context.Context, query string, args ...any) ([]Session, error) {
	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		var session Session
		if err := rows.Scan(
			&session.Id,
			&session.UserId,
			&session.ExpiresAt,
			&session.TwoFactorVerified,
		); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func (s *AuthRepository) fetchSession(ctx context.Context, query string, args ...any) (*Session, error) {
	var session Session
	err := s.db.QueryRow(ctx, query, args...).Scan(
		&session.Id,
		&session.UserId,
		&session.ExpiresAt,
		&session.TwoFactorVerified,
	)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// func (s *SessionRepository) updateSessionExpiration(ctx context.Context, session *Session) error {
// 	query := `
// 		UPDATE session
// 		SET expires_at = ?
// 		WHERE session.session_id = ?`

// 	return s.db.QueryRow(ctx, query, session.ExpiresAt, session.Id).Scan()
// }

// func (s *SessionRepository) invalidateSession(ctx context.Context, session *Session) error {
// 	query := `
// 		DELETE FROM session
// 		WHERE session.session_id = ?`

// 	return s.db.QueryRow(ctx, query, session.Id).Scan()
// }

// func (s *SessionRepository) Validate(ctx context.Context, sessionId string) (*SessionContext, error) {
// 	sessionId = hashSha256(sessionId)

// 	query := `
// 		SELECT *
// 		FROM session
// 		LEFT JOIN user ON session.user_id = user.user_id
// 		WHERE session_id = ?`

// 	sessionsContext, err := s.fetch(ctx, query, sessionId)

// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(sessionsContext) == 0 {
// 		return nil, errors.NotFound
// 	}

// 	sessionContext := sessionsContext[0]

// 	if time.Now().After(sessionContext.Session.ExpiresAt) {
// 		s.invalidateSession(ctx, sessionContext.Session)
// 		return nil, errors.SessionExpire
// 	}
// 	if time.Now().After(sessionContext.Session.ExpiresAt.Add(-sessionExpiresIn / 2)) {
// 		sessionContext.Session.ExpiresAt = time.Now().Add(sessionExpiresIn)
// 		s.updateSessionExpiration(ctx, sessionContext.Session)
// 	}

// 	return &sessionsContext[0], nil
// }

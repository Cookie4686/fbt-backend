package repo

import (
	"context"
	"fbt/backend/internal/domain/auth/model"

	"github.com/jackc/pgx/v5"
)

func (s *AuthRepository) CredentialsRegister(ctx context.Context, user *model.User, session *model.Session) error {
	batch := &pgx.Batch{}
	batch.Queue(`
		INSERT INTO users(user_id, username, email, email_verified, password, password_salt, password_enabled)
		VALUES (@userId, @username, @email, @emailVerified, @password, @passwordSalt, @passwordEnabled)
	`,
		pgx.NamedArgs{
			"userId":          user.Id,
			"username":        user.Username,
			"email":           user.Email,
			"emailVerified":   user.EmailVerified,
			"password":        user.Password,
			"passwordSalt":    user.PasswordSalt,
			"passwordEnabled": user.PasswordEnabled,
		},
	)
	batch.Queue(`
		INSERT INTO sessions(session_id, user_id, expires_at, two_factor_verified)
		VALUES (@sessionId, @sessionUserId, @expiresAt, @twoFactorVerified);
	`,
		pgx.NamedArgs{
			"sessionId":         session.Id,
			"sessionUserId":     session.UserId,
			"expiresAt":         session.ExpiresAt,
			"twoFactorVerified": session.TwoFactorVerified,
		},
	)

	br := s.db.SendBatch(ctx, batch)
	_, err := br.Exec()
	return err
}

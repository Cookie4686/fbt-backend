package credentials

import (
	"context"
	"fbt/backend/internal/domain/auth/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo interface {
	Register(context.Context, *model.User, *model.Session) error
}

type repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) Repo {
	return Repo(&repo{db: db})
}

func (s repo) Register(ctx context.Context, user *model.User, session *model.Session) error {
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
		INSERT INTO sessions(session_id, user_id, created_at, expires_at, two_factor_verified)
		VALUES (@sessionId, @sessionUserId, @createdAt, @expiresAt, @twoFactorVerified);
	`,
		pgx.NamedArgs{
			"sessionId":         session.Id,
			"sessionUserId":     session.UserId,
			"createdAt":         session.CreatedAt,
			"expiresAt":         session.ExpiresAt,
			"twoFactorVerified": session.TwoFactorVerified,
		},
	)

	_, err := s.db.SendBatch(ctx, batch).Exec()
	return err
}

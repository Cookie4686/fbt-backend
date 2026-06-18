package oauth

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

func (s *repo) GetUserOAuth(ctx context.Context, provider string, idToken string) (*model.UserOAuth, error) {
	query := `
		SELECT user_oauth.* FROM user_oauth
		LEFT JOIN oauth_providers ON user_oauth.oauth_provider_id = oauth_providers.oauth_provider_id
		WHERE oauth_providers.name = @provider AND user_oauth.id_token = @id_token
	`
	args := pgx.NamedArgs{"provider": provider, "id_token": idToken}
	userOAuth, err := s.fetchUserOAuth(ctx, query, args)
	if err != nil {
		if err == pgx.ErrNoRows {
			return userOAuth, errors.NotFound
		}
		return nil, err
	}
	return userOAuth, err
}

func (s *repo) LinkOAuth(ctx context.Context, provider string, userID string, idToken string) error {
	query := `
		INSERT INTO user_oauth(user_id, oauth_provider_id, id_token)
		VALUES (
			@user_id,
			(
				SELECT oauth_provider_id FROM oauth_providers
				WHERE oauth_providers.name = @provider
			),
			@id_token
		)
	`
	args := pgx.NamedArgs{"user_id": userID, "provider": provider, "id_token": idToken}
	_, err := s.db.Exec(ctx, query, args)
	return err
}

func (s *repo) CreateOAuthRegistration(ctx context.Context, provider string, oauthRegistration *model.OauthRegistration) error {
	query := `
		INSERT INTO oauth_registration(oauth_provider_id, id_token, registration_id, expires_at)
		VALUES (
			(
				SELECT oauth_provider_id FROM oauth_providers
				WHERE oauth_providers.name = @provider
			),
			@id_token,
			@registration_id,
			@expires_at
		)
	`
	args := pgx.NamedArgs{
		"provider":        provider,
		"id_token":        oauthRegistration.IDToken,
		"registration_id": oauthRegistration.RegistrationID,
		"expires_at":      oauthRegistration.ExpiresAt,
	}
	_, err := s.db.Exec(ctx, query, args)
	return err
}

func (s *repo) GetOAuthRegistration(ctx context.Context, registrationId string) (*model.OauthRegistration, error) {
	query := `
		SELECT * FROM oauth_registration
		WHERE registration_id = @registration_id
		ORDER BY expires_at DESC
	`
	args := pgx.NamedArgs{"registration_id": registrationId}
	OAuthRegistration, err := s.fetchOAuthRegistration(ctx, query, args)
	if err != nil {
		if err == pgx.ErrNoRows {
			return OAuthRegistration, errors.NotFound
		}
		return nil, err
	}
	return OAuthRegistration, nil
}

func (s *repo) DeleteOAuthRegistration(ctx context.Context, provider string, idToken string) error {
	query := `
		DELETE FROM oauth_registration
		WHERE oauth_provider_id = 
			(
				SELECT oauth_provider_id FROM oauth_providers
				WHERE oauth_providers.name = @provider
			)
		AND id_token = @idToken
	`
	args := pgx.NamedArgs{"provider": provider, "registration_id": idToken}
	_, err := s.db.Exec(ctx, query, args)
	return err
}

func (s *repo) OAuthRegister(ctx context.Context, registrationId string, user *model.User, session *model.Session) error {
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
		INSERT INTO user_oauth(oauth_provider_id, id_token, user_id)
		VALUES (
			(
				SELECT oauth_provider_id FROM oauth_registration
				WHERE registration_id = @registration_id
			),
			(
				SELECT id_token FROM oauth_registration
				WHERE registration_id = @registration_id
			),
			@user_id
		)
	`,
		pgx.NamedArgs{"user_id": user.Id, "registration_id": registrationId},
	)

	batch.Queue(`
		DELETE FROM oauth_registration
		WHERE registration_id = @registration_id
	`,
		pgx.NamedArgs{"registration_id": registrationId},
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

func (s *repo) GetUserProvider(ctx context.Context, userID string) ([]string, error) {
	query := `
		SELECT oauth_providers.name FROM user_oauth
		LEFT JOIN oauth_providers ON user_oauth.oauth_provider_id = oauth_providers.oauth_provider_id
		WHERE user_oauth.user_id = @user_id
	`
	args := pgx.NamedArgs{"user_id": userID}
	rows, err := s.db.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	var providers []string = make([]string, 0)
	for rows.Next() != false {
		var provider string
		err := rows.Scan(&provider)
		if err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}

	return providers, err
}

func (s *repo) fetchUserOAuths(ctx context.Context, query string, args ...any) ([]model.UserOAuth, error) {
	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userOAuths, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.UserOAuth])
	if err != nil {
		return nil, err
	}
	return userOAuths, nil
}

func (s *repo) fetchUserOAuth(ctx context.Context, query string, args ...any) (*model.UserOAuth, error) {
	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userOAuth, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserOAuth])
	if err != nil {
		return nil, err
	}
	return &userOAuth, nil
}

func (s *repo) fetchOAuthRegistration(ctx context.Context, query string, args ...any) (*model.OauthRegistration, error) {
	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	OAuthRegistration, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.OauthRegistration])
	if err != nil {
		return nil, err
	}
	return &OAuthRegistration, nil
}

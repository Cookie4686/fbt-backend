package webauthn

import (
	"context"
	"fbt/backend/internal/domain/auth/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo interface {
	GetPasskey(ctx context.Context, passkeyID string) (*model.Passkey, error)
	CreatePasskey(ctx context.Context, passkey *model.Passkey) error
	UpdatePasskeyCounter(ctx context.Context, passkeyID string, counter int64) (*model.Passkey, error)
}

type repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) Repo {
	return Repo(&repo{db})
}

func (s *repo) GetPasskey(ctx context.Context, passkeyID string) (*model.Passkey, error) {
	query := `
		SELECT * FROM user_passkey 
		WHERE passkey_id = @passkey_id
	`
	args := pgx.NamedArgs{"passkey_id": passkeyID}
	passkey, err := s.fetchPasskey(ctx, query, args)
	if err != nil {
		return nil, err
	}
	return passkey, nil
}

func (s *repo) CreatePasskey(ctx context.Context, passkey *model.Passkey) error {
	query := `
		INSERT INTO user_passkey(passkey_id, public_key, user_id, webauthn_user_id, counter, device_type, backed_up, transports)
		VALUES (@passkey_id, @public_key, @user_id, @webauthn_user_id, @counter, @device_type, @backed_up, @transports)
	`
	args := pgx.NamedArgs{
		`passkey_id`:       passkey.PasskeyID,
		`public_key`:       passkey.PublicKey,
		`user_id`:          passkey.UserID,
		"webauthn_user_id": passkey.WebauthnUserID,
		"counter":          passkey.Counter,
		"device_type":      passkey.DeviceType,
		"backed_up":        passkey.BackedUp,
		"transports":       passkey.Transports,
	}
	_, err := s.db.Exec(ctx, query, args)
	return err
}

func (s *repo) UpdatePasskeyCounter(ctx context.Context, passkeyID string, counter int64) (*model.Passkey, error) {
	query := `
		UPDATE user_passkey 
		SET counter = @counter
		WHERE passkey_id = @passkey_id
		RETURNING *
	`
	args := pgx.NamedArgs{"passkey_id": passkeyID, "counter": counter}
	passkey, err := s.fetchPasskey(ctx, query, args)
	if err != nil {
		return nil, err
	}
	return passkey, nil
}

func (s *repo) fetchPasskey(ctx context.Context, query string, args ...any) (*model.Passkey, error) {
	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	passkey, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Passkey])
	if err != nil {
		return nil, err
	}
	return &passkey, nil
}

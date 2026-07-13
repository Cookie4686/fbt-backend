package webauthn

import (
	"context"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/util"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo interface {
	GetPasskey(ctx context.Context, rpID string, credentialID string) (*model.WebAuthnCredential, error)
	CreatePasskey(ctx context.Context, credential *model.WebAuthnCredential) error
	UpdatePasskeyCounter(ctx context.Context, rpID string, credentialID string, counter int64) (*model.WebAuthnCredential, error)
}

type repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) Repo {
	return Repo(&repo{db})
}

func (s *repo) GetPasskey(ctx context.Context, rpID string, credentialID string) (*model.WebAuthnCredential, error) {
	query := `
		SELECT * FROM webauthn_credentials
		WHERE rp_id = @rp_id
		AND credential_id = @credential_id
	`
	args := pgx.StrictNamedArgs{"rp_id": rpID, "credential_id": credentialID}
	credential, err := util.FetchOne[model.WebAuthnCredential](s.db, ctx, query, args)
	if err != nil {
		return nil, err
	}
	return credential, nil
}

func (s *repo) CreatePasskey(ctx context.Context, credential *model.WebAuthnCredential) error {
	query := `
		INSERT INTO webauthn_credentials(rp_id, user_id, credential_id, public_key, counter, aaguid, device_type, transports, backup_state)
		VALUES (@rp_id, @user_id, @credential_id, @public_key, @counter, @aaguid, @device_type, @transports, @backup_state)
	`
	args := pgx.StrictNamedArgs{
		`rp_id`:         credential.RpID,
		`user_id`:       credential.UserID,
		`credential_id`: credential.CredentialID,
		`public_key`:    credential.PublicKey,
		"counter":       credential.Counter,
		`aaguid`:        credential.Aaguid,
		"device_type":   credential.DeviceType,
		"transports":    credential.Transports,
		"backup_state":  credential.BackupState,
	}
	_, err := s.db.Exec(ctx, query, args)
	return err
}

func (s *repo) UpdatePasskeyCounter(ctx context.Context, rpID string, credentialID string, counter int64) (*model.WebAuthnCredential, error) {
	query := `
		UPDATE webauthn_credentials 
		SET counter = @counter
		WHERE rp_id = @rp_id
		AND credential_id = @credential_id
		RETURNING *
	`
	args := pgx.NamedArgs{"rp_id": rpID, "credential_id": credentialID, "counter": counter}
	passkey, err := util.FetchOne[model.WebAuthnCredential](s.db, ctx, query, args)
	if err != nil {
		return nil, err
	}
	return passkey, nil
}

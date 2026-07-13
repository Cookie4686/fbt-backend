package webauthn_test

import (
	authv1 "fbt/backend/gen/proto/go/auth/v1"
	"fbt/backend/gen/proto/go/auth/v1/authv1connect"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/interceptor"
	"fbt/backend/internal/test"
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
)

func TestWebAuthn(t *testing.T) {
	ctx, baseURL := test.NewTestLocalAPI(t)

	client := authv1connect.NewWebAuthnServiceClient(http.DefaultClient, baseURL, connect.WithGRPC())

	session := test.SetupUser(t, ctx, baseURL)

	credential := &model.WebAuthnCredential{
		RpID:           "localhost",
		CredentialID:   "dasdasdasd",
		UserID:         session.UserId,
		PublicKey:      []byte{12, 23},
		Counter:        0,
		Aaguid:         []byte{12, 23},
		DeviceType:     "singleDevice",
		Transports:     []string{"nfc"},
		UserPresent:    false,
		UserVerified:   false,
		BackupEligible: false,
		BackupState:    false,
	}

	t.Run("Create User Passkey", func(t *testing.T) {
		ctx := interceptor.NewTokenContext(t.Context(), session.Id)
		res, err := client.CreateUserPasskey(ctx, &authv1.WebAuthnServiceCreateUserPasskeyRequest{
			Credential: &authv1.WebAuthnCredential{
				RpId:         credential.RpID,
				UserId:       credential.UserID,
				CredentialId: credential.CredentialID,
				PublicKey:    credential.PublicKey,
				Counter:      credential.Counter,
				Aaguid:       credential.Aaguid,
				DeviceType:   credential.DeviceType,
				Transports:   credential.Transports,
				BackedUp:     credential.BackupState,
			},
		},
		)
		require.NoError(t, err)
		require.NotNil(t, res.Credential)
	})

	t.Run("Get User Passkey", func(t *testing.T) {
		ctx := t.Context()
		res, err := client.GetUserPasskey(ctx, &authv1.WebAuthnServiceGetUserPasskeyRequest{
			RpId:         credential.RpID,
			CredentialId: credential.CredentialID,
		},
		)
		require.NoError(t, err)
		require.NotNil(t, res.Credential)
	})

	t.Run("Update User Passkey Counter", func(t *testing.T) {
		ctx := interceptor.NewTokenContext(t.Context(), session.Id)
		res, err := client.UpdatePasskeyCounter(ctx, &authv1.WebAuthnServiceUpdatePasskeyCounterRequest{
			RpId:         credential.RpID,
			CredentialId: credential.CredentialID,
			Counter:      credential.Counter + 1,
		},
		)
		require.NoError(t, err)
		require.NotNil(t, res.Session)
	})
}

package webauthn_test

import (
	"testing"

	"fbt/backend/internal/domain/auth/features/webauthn"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/test"

	authv1 "fbt/backend/gen/proto/go/auth/v1"

	"github.com/stretchr/testify/require"
)

func TestWebAuthn(t *testing.T) {
	d := test.Setup(t)
	testUtil := test.NewTestUtil(d)

	controller := webauthn.NewController(testUtil.AuthService, webauthn.NewRepo(d.DB))

	session := testUtil.SetupUser(t)

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
		ctx := testUtil.NewAuthContext(t.Context(), session.Id)
		res, err := controller.CreateUserPasskey(ctx, &authv1.WebAuthnServiceCreateUserPasskeyRequest{
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
		res, err := controller.GetUserPasskey(ctx, &authv1.WebAuthnServiceGetUserPasskeyRequest{
			RpId:         credential.RpID,
			CredentialId: credential.CredentialID,
		},
		)
		require.NoError(t, err)
		require.NotNil(t, res.Credential)
	})

	t.Run("Update User Passkey Counter", func(t *testing.T) {
		ctx := testUtil.NewAuthContext(t.Context(), session.Id)
		res, err := controller.UpdatePasskeyCounter(ctx, &authv1.WebAuthnServiceUpdatePasskeyCounterRequest{
			RpId:         credential.RpID,
			CredentialId: credential.CredentialID,
			Counter:      credential.Counter + 1,
		},
		)
		require.NoError(t, err)
		require.NotNil(t, res.Session)
	})
}

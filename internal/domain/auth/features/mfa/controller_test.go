package mfa_test

import (
	"testing"
	"time"

	authv1 "fbt/backend/gen/proto/go/auth/v1"
	"fbt/backend/internal/domain/auth/features/mfa"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/test"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMFA(t *testing.T) {
	d := test.Setup(t)
	service := service.NewService(d)

	session := test.SetupUser(t, d, service)
	client := mfa.NewController(service, mfa.NewRepo(d.DB))

	t.Run("Status", func(t *testing.T) {
		ctx := test.NewAuthContext(t.Context(), session.Id, service)

		res, err := client.Status(ctx, &authv1.MFAServiceStatusRequest{})
		require.NoError(t, err)

		assert.Equal(t, false, res.TotpEnabled)
	})

	key, err := totp.Generate(totp.GenerateOpts{Issuer: "fbt", AccountName: "test"})
	require.NoError(t, err)

	t.Run("TOTP Upsert", func(t *testing.T) {
		ctx := test.NewAuthContext(t.Context(), session.Id, service)

		res, err := client.TOTPUpsertKey(ctx, &authv1.MFAServiceTOTPUpsertKeyRequest{
			Key: key.Secret(),
		})
		require.NoError(t, err)

		assert.Equal(t, true, res.Session.TwoFactorVerified)

		session = res.Session
	})

	t.Run("TOTP Validate", func(t *testing.T) {
		ctx := test.NewAuthContext(t.Context(), session.Id, service)

		code, err := totp.GenerateCode(key.Secret(), time.Now())
		require.NoError(t, err)

		res, err := client.TOTPValidate(ctx, &authv1.MFAServiceTOTPValidateRequest{
			Code: code,
		})
		require.NoError(t, err)

		assert.Equal(t, true, res.Session.TwoFactorVerified)

		session = res.Session
	})

	t.Run("Status", func(t *testing.T) {
		ctx := test.NewAuthContext(t.Context(), session.Id, service)

		res, err := client.Status(ctx, &authv1.MFAServiceStatusRequest{})
		require.NoError(t, err)

		assert.Equal(t, true, res.TotpEnabled)
	})
}

package mfa_test

import (
	authv1 "fbt/backend/gen/proto/go/auth/v1"
	"fbt/backend/gen/proto/go/auth/v1/authv1connect"
	"fbt/backend/internal/interceptor"
	"fbt/backend/internal/test"
	"net/http"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMFA(t *testing.T) {
	ctx, baseURL := test.NewTestLocalAPI(t)

	client := authv1connect.NewMFAServiceClient(http.DefaultClient, baseURL, connect.WithGRPC())

	session := test.SetupUser(t, ctx, baseURL)

	t.Run("Status", func(t *testing.T) {
		ctx := interceptor.NewTokenContext(t.Context(), session.Id)

		res, err := client.Status(ctx, &authv1.MFAServiceStatusRequest{})
		require.NoError(t, err)

		assert.Equal(t, false, res.TotpEnabled)
	})

	key, err := totp.Generate(totp.GenerateOpts{Issuer: "fbt", AccountName: "test"})
	require.NoError(t, err)

	t.Run("TOTP Upsert", func(t *testing.T) {
		ctx := interceptor.NewTokenContext(t.Context(), session.Id)

		res, err := client.TOTPUpsertKey(ctx, &authv1.MFAServiceTOTPUpsertKeyRequest{
			Key: key.Secret(),
		})
		require.NoError(t, err)

		assert.Equal(t, true, res.Session.TwoFactorVerified)

		session = res.Session
	})

	t.Run("TOTP Validate", func(t *testing.T) {
		ctx := interceptor.NewTokenContext(t.Context(), session.Id)

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
		ctx := interceptor.NewTokenContext(t.Context(), session.Id)

		res, err := client.Status(ctx, &authv1.MFAServiceStatusRequest{})
		require.NoError(t, err)

		assert.Equal(t, true, res.TotpEnabled)
	})
}

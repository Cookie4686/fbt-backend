package mfa_test

import (
	"fbt/backend/internal/domain/auth/features/mfa/pb"
	"fbt/backend/internal/test"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func TestMFA(t *testing.T) {
	ctx, conn := test.NewTestLocalAPI(t)

	session := test.SetupUser(t, ctx, conn)

	client := pb.NewMFAClient(conn)

	t.Run("Status", func(t *testing.T) {
		ctx := metadata.AppendToOutgoingContext(t.Context(), "session_id", session.Id)

		res, err := client.Status(ctx, &pb.StatusRequest{})
		require.NoError(t, err)

		assert.Equal(t, false, res.TotpEnabled)
	})

	key, err := totp.Generate(totp.GenerateOpts{Issuer: "fbt", AccountName: "test"})
	require.NoError(t, err)

	t.Run("TOTP Upsert", func(t *testing.T) {
		ctx := metadata.AppendToOutgoingContext(t.Context(), "session_id", session.Id)

		res, err := client.TOTPUpsertKey(ctx, &pb.TOTPUpsertRequest{
			Key: key.Secret(),
		})
		require.NoError(t, err)

		assert.Equal(t, true, res.Session.TwoFactorVerified)

		session = res.Session
	})

	t.Run("TOTP Validate", func(t *testing.T) {
		ctx := metadata.AppendToOutgoingContext(t.Context(), "session_id", session.Id)

		code, err := totp.GenerateCode(key.Secret(), time.Now())
		require.NoError(t, err)

		res, err := client.TOTPValidate(ctx, &pb.TOTPValidateRequest{
			Code: code,
		})
		require.NoError(t, err)

		assert.Equal(t, true, res.IsValid)
	})

	t.Run("Status", func(t *testing.T) {
		ctx := metadata.AppendToOutgoingContext(t.Context(), "session_id", session.Id)

		res, err := client.Status(ctx, &pb.StatusRequest{})
		require.NoError(t, err)

		assert.Equal(t, true, res.TotpEnabled)
	})
}

package oauth_test

import (
	"fbt/backend/internal/domain/auth/common"
	"fbt/backend/internal/domain/auth/features/oauth/pb"
	"fbt/backend/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func TestOAuth(t *testing.T) {
	ctx, _, conn := test.NewTestConnection(t, 1234)

	client := pb.NewOAuthClient(conn)

	var registrationID string

	t.Run("Login", func(t *testing.T) {
		res, err := client.Login(ctx, &pb.LoginRequest{
			Token:    "token",
			Email:    "test@email.com",
			Provider: "google",
		})
		require.NoError(t, err)

		require.Equal(t, true, res.RegistrationNeeded)

		registrationID = res.RegistrationID
	})

	var session *common.Session
	t.Run("Register", func(t *testing.T) {
		res, err := client.Register(ctx, &pb.RegisterRequest{
			Username:        "test",
			Email:           "test@email.com",
			Password:        "12345678",
			PasswordEnabled: true,

			RegistrationID: registrationID,
			TokenID:        "token",
			Provider:       "google",
		})
		require.NoError(t, err)

		assert.Equal(t, false, res.Session.TwoFactorVerified)

		session = res.Session
	})

	t.Run("Status", func(t *testing.T) {
		ctx := metadata.AppendToOutgoingContext(t.Context(), "session_id", session.Id)

		res, err := client.Status(ctx, &pb.StatusRequest{})
		require.NoError(t, err)

		assert.Len(t, res.Providers, 1)
		assert.Contains(t, res.Providers, "google")
	})
}

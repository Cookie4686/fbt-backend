package credentials_test

import (
	authv1 "fbt/backend/gen/proto/go/auth/v1"
	"fbt/backend/gen/proto/go/auth/v1/authv1connect"
	"fbt/backend/internal/test"
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
)

func TestCredentials(t *testing.T) {
	ctx, baseURL := test.NewTestLocalAPI(t)

	client := authv1connect.NewCredentialServiceClient(http.DefaultClient, baseURL, connect.WithGRPC())

	username := "test"
	password := "12345678"

	t.Run("Register", func(t *testing.T) {
		ctx := t.Context()

		res, err := client.Register(ctx, &authv1.CredentialServiceRegisterRequest{
			Username: username,
			Email:    "test@email.com",
			Password: password,
		})
		require.NoError(t, err)

		require.Equal(t, false, res.Session.TwoFactorVerified)
	})

	t.Run("Login", func(t *testing.T) {
		res, err := client.Login(ctx, &authv1.CredentialServiceLoginRequest{
			Username: username,
			Password: password,
		})
		require.NoError(t, err)

		require.Equal(t, false, res.Session.TwoFactorVerified)
	})
}

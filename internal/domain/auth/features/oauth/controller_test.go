package oauth_test

import (
	authv1 "fbt/backend/gen/proto/go/auth/v1"
	"fbt/backend/gen/proto/go/auth/v1/authv1connect"
	"fbt/backend/internal/interceptor"
	"fbt/backend/internal/test"
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOAuth(t *testing.T) {
	_, baseURL := test.NewTestLocalAPI(t)

	client := authv1connect.NewOAuthServiceClient(http.DefaultClient, baseURL, connect.WithGRPC())

	var registrationID string

	email := "test@email.com"

	t.Run("Login", func(t *testing.T) {
		res, err := client.Login(t.Context(), &authv1.OAuthServiceLoginRequest{
			Token:    "token",
			Provider: "google",
			Email:    &email,
		})
		require.NoError(t, err)

		require.Equal(t, true, res.RegistrationNeeded)

		registrationID = res.RegistrationId
	})

	var session *authv1.Session

	t.Run("Register", func(t *testing.T) {
		res, err := client.Register(t.Context(), &authv1.OAuthServiceRegisterRequest{
			Username:        "test",
			Email:           "test@email.com",
			Password:        "12345678",
			PasswordEnabled: true,

			RegistrationId: registrationID,
			TokenId:        "token",
			Provider:       "google",
		})
		require.NoError(t, err)

		assert.Equal(t, false, res.Session.TwoFactorVerified)

		session = res.Session
	})

	t.Run("Status", func(t *testing.T) {
		ctx := interceptor.NewTokenContext(t.Context(), session.Id)

		res, err := client.Status(ctx, &authv1.OAuthServiceStatusRequest{})
		require.NoError(t, err)

		assert.Len(t, res.Providers, 1)
		assert.Contains(t, res.Providers, "google")
	})
}

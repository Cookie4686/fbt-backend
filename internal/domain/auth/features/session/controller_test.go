package session_test

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

func TestSession(t *testing.T) {
	ctx, baseURL := test.NewTestLocalAPI(t)

	client := authv1connect.NewSessionServiceClient(http.DefaultClient, baseURL, connect.WithGRPC())

	session := test.SetupUser(t, ctx, baseURL)

	t.Run("Validate", func(t *testing.T) {
		ctx := interceptor.NewTokenContext(t.Context(), session.Id)

		res, err := client.Validate(ctx, &authv1.SessionServiceValidateRequest{})
		require.NoError(t, err)

		assert.Equal(t, session.Id, res.Session.Id)
		assert.Equal(t, session.UserId, res.User.Id)
	})

	t.Run("Logout", func(t *testing.T) {
		ctx := interceptor.NewTokenContext(t.Context(), session.Id)

		_, err := client.Logout(ctx, &authv1.SessionServiceLogoutRequest{})
		require.NoError(t, err)
	})

	t.Run("Validate", func(t *testing.T) {
		ctx := interceptor.NewTokenContext(t.Context(), session.Id)

		_, err := client.Validate(ctx, &authv1.SessionServiceValidateRequest{})
		require.Error(t, err)
	})
}

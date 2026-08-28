package session_test

import (
	"testing"

	"fbt/backend/internal/domain/auth/features/session"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/test"

	authv1 "fbt/backend/gen/proto/go/auth/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSession(t *testing.T) {
	d := test.Setup(t)

	service := service.NewService(d)
	controller := session.NewController(service)
	session := test.SetupUser(t, d, nil)

	t.Run("Validate", func(t *testing.T) {
		ctx := test.NewAuthContext(t.Context(), session.Id, service)

		res, err := controller.Validate(ctx, &authv1.SessionServiceValidateRequest{})
		require.NoError(t, err)

		assert.Equal(t, session.Id, res.Session.Id)
		assert.Equal(t, session.UserId, res.User.Id)
	})

	t.Run("Logout", func(t *testing.T) {
		ctx := test.NewAuthContext(t.Context(), session.Id, service)

		_, err := controller.Logout(ctx, &authv1.SessionServiceLogoutRequest{})
		require.NoError(t, err)
	})

	t.Run("Validate", func(t *testing.T) {
		ctx := test.NewAuthContext(t.Context(), session.Id, service)

		_, err := controller.Validate(ctx, &authv1.SessionServiceValidateRequest{})
		require.Error(t, err)
	})
}

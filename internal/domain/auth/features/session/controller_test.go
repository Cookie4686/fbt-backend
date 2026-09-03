package session_test

import (
	"testing"

	"fbt/backend/internal/domain/auth/features/session"
	"fbt/backend/internal/test"

	authv1 "fbt/backend/gen/proto/go/auth/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSession(t *testing.T) {
	d := test.Setup(t)
	testUtil := test.NewTestUtil(d)

	controller := session.NewController(testUtil.AuthService)

	session := testUtil.SetupUser(t)

	t.Run("Validate", func(t *testing.T) {
		ctx := testUtil.NewAuthContext(t.Context(), session.Id)

		res, err := controller.Validate(ctx, &authv1.SessionServiceValidateRequest{})
		require.NoError(t, err)

		assert.Equal(t, session.Id, res.Session.Id)
		assert.Equal(t, session.UserId, res.User.Id)
	})

	t.Run("Logout", func(t *testing.T) {
		ctx := testUtil.NewAuthContext(t.Context(), session.Id)

		_, err := controller.Logout(ctx, &authv1.SessionServiceLogoutRequest{})
		require.NoError(t, err)
	})

	t.Run("Validate", func(t *testing.T) {
		ctx := testUtil.NewAuthContext(t.Context(), session.Id)

		_, err := controller.Validate(ctx, &authv1.SessionServiceValidateRequest{})
		require.Error(t, err)
	})
}

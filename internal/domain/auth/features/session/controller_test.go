package session_test

import (
	"fbt/backend/internal/domain/auth/features/session/pb"
	"fbt/backend/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func TestSession(t *testing.T) {
	ctx, conn := test.NewTestLocalAPI(t)

	session := test.SetupUser(t, ctx, conn)

	client := pb.NewSessionClient(conn)

	t.Run("Validate", func(t *testing.T) {
		ctx := metadata.AppendToOutgoingContext(t.Context(), "session_id", session.Id)

		res, err := client.Validate(ctx, &pb.ValidateRequest{})
		require.NoError(t, err)

		assert.Equal(t, session.Id, res.Session.Id)
		assert.Equal(t, session.UserID, res.User.Id)
	})

	t.Run("Logout", func(t *testing.T) {
		ctx := metadata.AppendToOutgoingContext(t.Context(), "session_id", session.Id)

		_, err := client.Logout(ctx, &pb.LogoutRequest{})
		require.NoError(t, err)
	})

	t.Run("Validate", func(t *testing.T) {
		ctx := metadata.AppendToOutgoingContext(t.Context(), "session_id", session.Id)

		_, err := client.Validate(ctx, &pb.ValidateRequest{})
		require.Error(t, err)
	})
}

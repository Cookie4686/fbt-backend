package test

import (
	"context"
	"fbt/backend/internal/domain/auth/common"
	"fbt/backend/internal/domain/auth/features/credentials/pb"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func SetupUser(t *testing.T, ctx context.Context, conn *grpc.ClientConn) *common.Session {
	client := pb.NewCredentialsClient(conn)

	res, err := client.Register(ctx, &pb.RegisterRequest{
		Username: "test",
		Email:    "test@email.com",
		Password: "12345678",
	})
	require.NoError(t, err)

	return res.Session
}

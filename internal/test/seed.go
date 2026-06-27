package test

import (
	"context"
	authv1 "fbt/backend/gen/proto/go/auth/v1"
	"fbt/backend/gen/proto/go/auth/v1/authv1connect"
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
)

func SetupUser(t *testing.T, ctx context.Context, baseURL string) *authv1.Session {
	client := authv1connect.NewCredentialServiceClient(http.DefaultClient, baseURL, connect.WithGRPC())

	res, err := client.Register(ctx, &authv1.CredentialServiceRegisterRequest{
		Username: "test",
		Email:    "test@email.com",
		Password: "12345678",
	})
	require.NoError(t, err)

	return res.Session
}

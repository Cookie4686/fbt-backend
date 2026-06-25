package test

import (
	"context"
	"fbt/backend/internal/server"
	"fbt/backend/internal/util"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewTestLocalAPI(t *testing.T) (context.Context, *grpc.ClientConn) {
	ctx := t.Context()

	ChangeDirectory(t)

	d, err := util.NewDependency(ctx, ".env.test")
	require.NoError(t, err)

	ClearDatabase(t, ctx, d.DB)

	s := server.NewServer(d)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", d.CFG.API.PORT))
	require.NoError(t, err)

	go func() {
		err := s.Serve(lis)
		require.NoError(t, err)
		t.Cleanup((func() {
			s.Stop()
		}))
	}()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	t.Cleanup(func() {
		conn.Close()
	})
	require.NoError(t, err)

	return ctx, conn
}

func ClearDatabase(t *testing.T, ctx context.Context, db *pgxpool.Pool) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	for _, v := range []string{
		"/sqlc/schema/auth_down.sql",
		"/sqlc/schema/bookkeeping_down.sql",
		"/sqlc/schema/auth_up.sql",
		"/sqlc/schema/bookkeeping_up.sql",
		"/sqlc/seed/auth.sql",
	} {
		file, err := os.ReadFile(filepath.Join(wd, v))
		require.NoError(t, err)

		_, err = db.Exec(ctx, string(file))
		require.NoError(t, err)
	}
}

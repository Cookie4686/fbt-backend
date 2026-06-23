package test

import (
	"context"
	"fbt/backend/internal/api"
	"fbt/backend/internal/config"
	"fbt/backend/internal/util"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func NewTestContainer(t *testing.T, dbName string, user string, password string) (context.Context, *postgres.PostgresContainer, *pgxpool.Pool) {
	ChangeDirectory(t)

	ctx := context.Background()

	ctr, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(user),
		postgres.WithPassword(password),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)

	testcontainers.CleanupContainer(t, ctr)
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)

	for _, v := range []string{
		"/sqlc/schema/auth_up.sql",
		"/sqlc/schema/bookkeeping_up.sql",
	} {
		migrationFile, err := os.ReadFile(filepath.Join(wd, v))
		require.NoError(t, err)

		_, _, err = ctr.Exec(ctx, []string{"psql", "-U", user, "-d", dbName, "-c", string(migrationFile)})
		require.NoError(t, err)
	}

	for _, v := range []string{
		"/sqlc/seed/auth.sql",
	} {
		seedFile, err := os.ReadFile(filepath.Join(wd, v))
		require.NoError(t, err)

		_, _, err = ctr.Exec(ctx, []string{"psql", "-U", user, "-d", dbName, "-c", string(seedFile)})
		require.NoError(t, err)
	}

	err = ctr.Snapshot(ctx)
	require.NoError(t, err)

	conn, err := ctr.ConnectionString(ctx)
	require.NoError(t, err)

	db, err := pgxpool.New(ctx, conn)
	require.NoError(t, err)

	return ctx, ctr, db
}

func NewTestAPI(t *testing.T) (context.Context, *postgres.PostgresContainer, *grpc.Server) {
	name := "test"
	user := "user"
	password := "password"

	ctx, ctr, db := NewTestContainer(t, name, user, password)

	cfg, err := config.LoadConfig()
	require.NoError(t, err)

	logger, err := util.NewLogger(cfg)
	require.NoError(t, err)

	svr := api.NewGRPCServer(logger, db, cfg)

	return ctx, ctr, svr
}

func NewTestConnection(t *testing.T, port int) (context.Context, *postgres.PostgresContainer, *grpc.ClientConn) {
	ctx, ctr, svr := NewTestAPI(t)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	require.NoError(t, err)

	go func() {
		if err := svr.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	return ctx, ctr, conn
}

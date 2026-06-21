package test

import (
	"context"
	"fbt/backend/internal/api"
	"fbt/backend/internal/config"
	"fbt/backend/internal/util"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

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

	err = ctr.Snapshot(ctx)
	require.NoError(t, err)

	conn, err := ctr.ConnectionString(ctx)
	require.NoError(t, err)

	db, err := pgxpool.New(ctx, conn)
	require.NoError(t, err)

	return ctx, ctr, db
}

func NewTestAPI(t *testing.T) (context.Context, *postgres.PostgresContainer, *httptest.Server) {
	name := "test"
	user := "user"
	password := "password"

	ctx, ctr, db := NewTestContainer(t, name, user, password)

	cfg, err := config.LoadConfig()
	require.NoError(t, err)

	logger, err := util.NewLogger(cfg)
	require.NoError(t, err)

	handler := api.NewAPIHandler(logger, db, cfg)
	svr := httptest.NewServer(handler)

	return ctx, ctr, svr
}

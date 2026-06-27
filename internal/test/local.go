package test

import (
	"context"
	"fbt/backend/internal/server"
	"fbt/backend/internal/util"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func NewTestLocalAPI(t *testing.T) (ctx context.Context, baseURL string) {
	ctx = t.Context()

	ChangeDirectory(t)

	d, err := util.NewDependency(ctx, ".env.test")
	require.NoError(t, err)

	ClearDatabase(t, ctx, d.DB)

	svr := server.NewServer(d)

	go func() { server.StartListening(svr, d) }()

	return ctx, fmt.Sprintf("http://localhost:%d", d.CFG.API.PORT)
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

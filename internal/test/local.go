package test

import (
	"context"
	"database/sql"
	"fbt/backend/internal/server"
	"fbt/backend/internal/util"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
)

func NewTestLocalAPI(t *testing.T) (ctx context.Context, baseURL string) {
	ctx = t.Context()

	ChangeDirectory(t)

	d, err := util.NewDependency(ctx)
	require.NoError(t, err)

	ClearDatabase(t, ctx, d.DB)

	svr := server.NewServer(d)

	go func() { server.StartListening(svr, d) }()

	return ctx, fmt.Sprintf("http://localhost:%d", d.CFG.API.PORT)
}

func ClearDatabase(t *testing.T, ctx context.Context, db *pgxpool.Pool) {
	sqlDB, err := sql.Open("pgx", db.Config().ConnString())
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	gooseProvider, err := goose.NewProvider("postgres", sqlDB, os.DirFS(filepath.Join(wd, "/sqlc/schema")))
	require.NoError(t, err)

	current, _, err := gooseProvider.GetVersions(ctx)
	require.NoError(t, err)

	if current != 0 {
		_, err = gooseProvider.DownTo(ctx, 0)
		require.NoError(t, err)
	}

	_, err = gooseProvider.Up(ctx)
	require.NoError(t, err)
}

// Package test for testing utilities
package test

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"fbt/backend/internal/util"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
)

func SetupDatabase(t *testing.T, cfg *util.Config) *pgxpool.Pool {
	ctx := t.Context()

	ClearAndMigrateDatabase(t, cfg.PGURL)

	db, err := util.NewDatabasePool(ctx, cfg.PGURL)
	require.NoError(t, err)

	return db
}

func ClearAndMigrateDatabase(t *testing.T, pgurl string) {
	ctx := t.Context()

	sqlDB, err := sql.Open("pgx", pgurl)
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

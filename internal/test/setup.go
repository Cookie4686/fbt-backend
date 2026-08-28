// Package test for testing utilities
package test

import (
	"os"
	"testing"

	"fbt/backend/internal/util"

	"github.com/stretchr/testify/require"
)

func Setup(t *testing.T) *util.Dependency {
	ChangeDirectory(t)

	cfg, err := util.NewConfig()
	require.NoError(t, err)

	db := SetupDatabase(t, cfg)

	logger, err := util.NewLogger(cfg)
	require.NoError(t, err)

	mail, err := util.NewMailClient(cfg)
	require.NoError(t, err)

	return &util.Dependency{Logger: logger, DB: db, CFG: cfg, Mail: mail}
}

func ChangeDirectory(t *testing.T) {
	var cwd string

	cwd, ok := os.LookupEnv("CWD")
	if !ok {
		if cwd, ok = os.LookupEnv("cwd"); !ok {
			t.Fatal("No CWD environment defined")
		}
	}

	if cwd == "" {
		t.Fatal("CWD environment cannot be empty")
	}

	if err := os.Chdir(cwd); err != nil {
		t.Fatal("Chdir error:", err)
	}
}

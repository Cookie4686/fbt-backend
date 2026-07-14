package util

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func FetchOne[T any](db *pgxpool.Pool, ctx context.Context, query string, args ...any) (*T, error) {
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[T])
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func FetchMany[T any](db *pgxpool.Pool, ctx context.Context, query string, args ...any) (*[]T, error) {
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[T])
	if err != nil {
		return nil, err
	}

	return &models, nil
}

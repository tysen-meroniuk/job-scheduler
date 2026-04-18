package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect opens a pgx connection pool against the given DSN.
func Connect(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	return pgxpool.NewWithConfig(ctx, cfg)
}

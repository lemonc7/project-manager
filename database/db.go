package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lemonc7/project-manager/config"
)

func NewDB(cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := dbPool.Ping(ctx); err != nil {
		dbPool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return dbPool, nil
}

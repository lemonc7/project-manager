package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lemonc7/project-manager/config"
)

func NewDB(cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	ctx := context.Background()
	dbPool, err := connectDB(ctx, cfg)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "3D000" {
			if err := createDatabase(cfg); err != nil {
				return nil, fmt.Errorf("create database: %w", err)
			}
			return connectDB(ctx, cfg)
		}
		return nil, err
	}

	return dbPool, nil
}

func connectDB(ctx context.Context, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	dbPool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	if err := dbPool.Ping(ctx); err != nil {
		dbPool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return dbPool, nil
}

func createDatabase(cfg config.DatabaseConfig) error {
	adminDSN := fmt.Sprintf("postgres://%s:%s@%s:%d/postgres?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.SSLMode)

	ctx := context.Background()
	adminPool, err := pgxpool.New(ctx, adminDSN)
	if err != nil {
		return fmt.Errorf("connect to admin database: %w", err)
	}
	defer adminPool.Close()

	safeName := pgx.Identifier{cfg.DBname}.Sanitize()
	_, err = adminPool.Exec(ctx, "CREATE DATABASE "+safeName)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "42P04" {
			return nil
		}
		return fmt.Errorf("create database: %w", err)
	}

	return nil
}

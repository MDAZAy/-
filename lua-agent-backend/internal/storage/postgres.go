package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	defaultMaxOpenConns    = 10
	defaultMaxIdleConns    = 5
	defaultConnMaxLifetime = 30 * time.Minute
	defaultConnMaxIdleTime = 5 * time.Minute
)

// PostgresConfig controls DB connectivity and pool settings.
type PostgresConfig struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// PostgresRepository is a PostgreSQL-backed implementation of Repository.
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository opens a connection pool and verifies it with PingContext.
func NewPostgresRepository(ctx context.Context, cfg PostgresConfig) (*PostgresRepository, error) {
	if cfg.DSN == "" {
		return nil, fmt.Errorf("postgres dsn is required")
	}

	db, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("open postgres connection: %w", err)
	}

	cfg = withDefaultPoolConfig(cfg)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return &PostgresRepository{db: db}, nil
}

// Close releases the underlying connection pool.
func (r *PostgresRepository) Close() error {
	if r == nil || r.db == nil {
		return nil
	}

	return r.db.Close()
}

// DB exposes the underlying sql.DB for migrations and health checks.
func (r *PostgresRepository) DB() *sql.DB {
	if r == nil {
		return nil
	}

	return r.db
}

func withDefaultPoolConfig(cfg PostgresConfig) PostgresConfig {
	if cfg.MaxOpenConns <= 0 {
		cfg.MaxOpenConns = defaultMaxOpenConns
	}
	if cfg.MaxIdleConns <= 0 {
		cfg.MaxIdleConns = defaultMaxIdleConns
	}
	if cfg.ConnMaxLifetime <= 0 {
		cfg.ConnMaxLifetime = defaultConnMaxLifetime
	}
	if cfg.ConnMaxIdleTime <= 0 {
		cfg.ConnMaxIdleTime = defaultConnMaxIdleTime
	}

	return cfg
}

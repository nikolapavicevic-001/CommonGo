// Package postgres provides pgxpool connection helpers and configuration.
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Config holds PostgreSQL connection configuration.
type Config struct {
	// URL is the PostgreSQL connection string (e.g., postgres://user:pass@host:5432/db)
	URL string

	// MaxConns is the maximum number of connections in the pool (default: 10)
	MaxConns int32

	// MinConns is the minimum number of connections in the pool (default: 2)
	MinConns int32

	// MaxConnLifetime is the maximum lifetime of a connection (default: 1 hour)
	MaxConnLifetime time.Duration

	// MaxConnIdleTime is the maximum idle time for a connection (default: 30 minutes)
	MaxConnIdleTime time.Duration

	// HealthCheckPeriod is how often to check connection health (default: 1 minute)
	HealthCheckPeriod time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig(url string) Config {
	return Config{
		URL:               url,
		MaxConns:          10,
		MinConns:          2,
		MaxConnLifetime:   time.Hour,
		MaxConnIdleTime:   30 * time.Minute,
		HealthCheckPeriod: time.Minute,
	}
}

// Open creates and returns a new pgxpool.Pool using the provided configuration.
// It validates the connection by pinging the database before returning.
func Open(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("parsing postgres config: %w", err)
	}

	// Apply pool settings
	if cfg.MaxConns > 0 {
		poolCfg.MaxConns = cfg.MaxConns
	}
	if cfg.MinConns > 0 {
		poolCfg.MinConns = cfg.MinConns
	}
	if cfg.MaxConnLifetime > 0 {
		poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	}
	if cfg.MaxConnIdleTime > 0 {
		poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	}
	if cfg.HealthCheckPeriod > 0 {
		poolCfg.HealthCheckPeriod = cfg.HealthCheckPeriod
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("creating postgres pool: %w", err)
	}

	// Validate connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pinging postgres: %w", err)
	}

	return pool, nil
}

// MustOpen is like Open but panics on error.
func MustOpen(ctx context.Context, cfg Config) *pgxpool.Pool {
	pool, err := Open(ctx, cfg)
	if err != nil {
		panic(err)
	}
	return pool
}


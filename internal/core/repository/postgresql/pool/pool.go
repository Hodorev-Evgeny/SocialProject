package core_repository_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool interface {
	QueryRow(ctx context.Context, sql string, args ...any) (pgx.Row, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	GetTimeout() time.Duration
	Close()
}

type ConnectionPool struct {
	pool    *pgxpool.Pool
	timeout time.Duration
}

func CreatePool(ctx context.Context, config PostgresConfig) (*ConnectionPool, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	parseconf, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres url: %w", err)
	}

	pgpool, err := pgxpool.NewWithConfig(ctx, parseconf)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	return &ConnectionPool{
		pool:    pgpool,
		timeout: config.Timeout,
	}, nil
}

func (c *ConnectionPool) GetTimeout() time.Duration {
	return c.timeout
}

func (c *ConnectionPool) Close() {
	c.pool.Close()
}

func (c *ConnectionPool) Ping(ctx context.Context) error {
	return c.pool.Ping(ctx)
}

func (c *ConnectionPool) QueryRow(ctx context.Context, sql string, args ...any) (pgx.Row, error) {
	return c.pool.QueryRow(ctx, sql, args...), nil
}

func (c *ConnectionPool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return c.pool.Query(ctx, sql, args...)
}

func CreatePoolMust(ctx context.Context, config PostgresConfig) *ConnectionPool {
	poolconnect, err := CreatePool(ctx, config)
	if err != nil {
		panic(err)
	}

	return poolconnect
}

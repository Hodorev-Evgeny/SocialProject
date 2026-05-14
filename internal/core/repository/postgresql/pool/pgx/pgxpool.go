package core_pgx_pool

import (
	"context"
	"fmt"
	"time"

	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	*pgxpool.Pool
	timeout time.Duration
}

func CreatePool(ctx context.Context, config PostgresConfig) (*Pool, error) {
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

	return &Pool{
		Pool:    pgpool,
		timeout: config.Timeout,
	}, nil
}

func (c *Pool) GetTimeout() time.Duration {
	return c.timeout
}

func (c *Pool) Exec(ctx context.Context, sql string, arguments ...any) (core_repository_pool.CommandTag, error) {
	ans, err := c.Pool.Exec(ctx, sql, arguments...)
	if err != nil {
		return nil, err
	}
	return CommandTag{ans}, nil
}

func (c *Pool) Ping(ctx context.Context) error {
	return c.Pool.Ping(ctx)
}

func (c *Pool) QueryRow(ctx context.Context, sql string, args ...any) core_repository_pool.Row {
	ans := c.Pool.QueryRow(ctx, sql, args...)
	return Row{ans}
}

func (c *Pool) Query(ctx context.Context, sql string, args ...any) (core_repository_pool.Rows, error) {
	ans, err := c.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return Rows{ans}, nil
}

func (c *Pool) Begin(ctx context.Context) (core_repository_pool.Tx, error) {
	t, err := c.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &Tx{tx: t}, nil
}

func CreatePoolMust(ctx context.Context, config PostgresConfig) *Pool {
	poolconnect, err := CreatePool(ctx, config)
	if err != nil {
		panic(err)
	}

	return poolconnect
}

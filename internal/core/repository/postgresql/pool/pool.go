package core_repository_pool

import (
	"context"
	"time"
)

type Pool interface {
	QueryRow(ctx context.Context, sql string, args ...any) Row
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	Exec(ctx context.Context, sql string, arguments ...any) (CommandTag, error)
	Begin(ctx context.Context) (Tx, error)
	GetTimeout() time.Duration
	Close()
}

type Tx interface {
	QueryRow(ctx context.Context, sql string, args ...any) Row
	Exec(ctx context.Context, sql string, arguments ...any) (CommandTag, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type Row interface {
	Scan(dest ...any) error
}

type Rows interface {
	Close()
	Err() error
	Next() bool
	Scan(dest ...any) error
}

type CommandTag interface {
	RowsAffected() int64
}

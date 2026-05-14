package core_pgx_pool

import (
	"context"

	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
	"github.com/jackc/pgx/v5"
)

type Tx struct {
	tx pgx.Tx
}

func (t *Tx) Exec(ctx context.Context, sql string, arguments ...any) (core_repository_pool.CommandTag, error) {
	ans, err := t.tx.Exec(ctx, sql, arguments...)
	if err != nil {
		return nil, err
	}
	return CommandTag{ans}, nil
}

func (t *Tx) QueryRow(ctx context.Context, sql string, args ...any) core_repository_pool.Row {
	return Row{t.tx.QueryRow(ctx, sql, args...)}
}

func (t *Tx) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *Tx) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

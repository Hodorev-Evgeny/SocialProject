package core_pgx_pool

import (
	"errors"

	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Row struct {
	pgx.Row
}

type Rows struct {
	pgx.Rows
}

func (r Rows) Scan(dest ...any) error {
	err := r.Rows.Scan(dest...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_repository_pool.ErrNoRows
		}
		return err
	}
	return nil
}

type CommandTag struct {
	pgconn.CommandTag
}

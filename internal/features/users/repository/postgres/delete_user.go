package features_users_repository

import (
	"context"
	"fmt"

	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
)

func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	query := `
		DELETE 
		FROM trackerapp.users 
		WHERE id = $1;`

	comantag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf(`user with id %d: %w`, id, err)
	}
	if comantag.RowsAffected() == 0 {
		return fmt.Errorf(`user with id %d: %w`, id, core_errors.ErrorNotFoud)
	}

	return nil
}

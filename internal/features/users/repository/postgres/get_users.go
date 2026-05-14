package features_users_repository

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
)

func (r *UserRepository) GetUsers(
	ctx context.Context,
	limit *int, offset *int,
) ([]core_domain.User, error) {
	ctx, close := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer close()

	query := `
		SELECT *
		FROM trackerapp.users
		ORDER BY id
		LIMIT $1 OFFSET $2;`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error get users: %w", err)
	}
	defer rows.Close()

	var users []core_domain.User
	for rows.Next() {
		var user core_domain.User

		err := rows.Scan(
			&user.ID,
			&user.Full_name,
			&user.Email,
			&user.Phone_number,
			&user.Password,
			&user.Time_add,
			&user.Role,
			&user.Is_verified)
		if err != nil {
			if errors.Is(err, core_repository_pool.ErrNoRows) {
				return nil, fmt.Errorf("user not in database: %w", core_errors.ErrorNotFoud)
			}
			return nil, fmt.Errorf("error scan users: %w", err)
		}

		users = append(users, user)
	}

	return users, nil
}

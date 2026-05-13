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
		SELECT id, full_name, email, phone_number, password, role, time_add
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
		var roleValue string

		err := rows.Scan(
			&user.ID,
			&user.Full_name,
			&user.Email,
			&user.Phone_number,
			&user.Password,
			&roleValue,
			&user.Time_add,
		)
		if err != nil {
			if errors.Is(err, core_repository_pool.ErrNoRows) {
				return nil, fmt.Errorf("user not in database: %w", core_errors.ErrorNotFoud)
			}
			return nil, fmt.Errorf("error scan users: %w", err)
		}

		user.Role = core_domain.UserRole(roleValue)
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error rows users: %w", err)
	}

	return users, nil
}

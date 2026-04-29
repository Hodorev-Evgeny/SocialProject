package features_users_repository

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (r *UserRepository) PatchUser(ctx context.Context,
	id int,
	patch core_domain.User,
) (core_domain.User, error) {
	query := `UPDATE trackerapp.users 
			SET full_name=$1, email=$2, phone_number=$3
			WHERE id=$4
			RETURNING id, full_name, email, phone_number, password, time_add;`

	row := r.pool.QueryRow(ctx,
		query,
		patch.Full_name,
		patch.Email,
		patch.Phone_number,
		id,
	)

	var UserUpdated core_domain.User
	err := row.Scan(
		&UserUpdated.ID,
		&UserUpdated.Full_name,
		&UserUpdated.Email,
		&UserUpdated.Phone_number,
		&UserUpdated.Password,
		&UserUpdated.Time_add)
	if err != nil {
		return core_domain.User{}, fmt.Errorf("failed to update user: %w", err)
	}

	return UserUpdated, nil
}

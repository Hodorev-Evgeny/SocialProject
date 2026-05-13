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
	query := `
		UPDATE trackerapp.users
		SET full_name=$1, email=$2, phone_number=$3
		WHERE id=$4
		RETURNING id, full_name, email, phone_number, password, role, time_add;`

	row := r.pool.QueryRow(ctx,
		query,
		patch.Full_name,
		patch.Email,
		patch.Phone_number,
		id,
	)

	var userUpdated core_domain.User
	var roleValue string
	err := row.Scan(
		&userUpdated.ID,
		&userUpdated.Full_name,
		&userUpdated.Email,
		&userUpdated.Phone_number,
		&userUpdated.Password,
		&roleValue,
		&userUpdated.Time_add,
	)
	if err != nil {
		return core_domain.User{}, fmt.Errorf("failed to update user: %w", err)
	}

	userUpdated.Role = core_domain.UserRole(roleValue)

	return userUpdated, nil
}

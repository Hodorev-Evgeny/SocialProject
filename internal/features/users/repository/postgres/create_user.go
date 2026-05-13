package features_users_repository

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (r *UserRepository) AddUser(ctx context.Context, user core_domain.User) (core_domain.User, error) {
	ctx, close := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer close()

	query := `
		INSERT INTO trackerapp.users (full_name, email, phone_number, password, role, time_add)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, full_name, email, phone_number, password, role, time_add;`

	row := r.pool.QueryRow(ctx, query,
		user.Full_name,
		user.Email,
		user.Phone_number,
		user.Password,
		string(user.Role),
		user.Time_add,
	)

	var userDomain core_domain.User
	var roleValue string
	err := row.Scan(
		&userDomain.ID,
		&userDomain.Full_name,
		&userDomain.Email,
		&userDomain.Phone_number,
		&userDomain.Password,
		&roleValue,
		&userDomain.Time_add,
	)
	if err != nil {
		return core_domain.User{}, fmt.Errorf("error inserting user: %w", err)
	}

	userDomain.Role = core_domain.UserRole(roleValue)

	return userDomain, nil
}

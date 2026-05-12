package features_users_repository

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (r *UserRepository) AddUser(ctx context.Context, user core_domain.User) (core_domain.User, error) {
	ctx, close := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer close()

	quqry := `
	INSERT INTO trackerapp.users (full_name, email, phone_number, password, time_add)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, full_name, email, phone_number, password, time_add, role, is_verified;
	`

	row := r.pool.QueryRow(ctx, quqry,
		user.Full_name,
		user.Email,
		user.Phone_number,
		user.Password,
		user.Time_add)

	var userDomain core_domain.User
	err := row.Scan(
		&userDomain.ID,
		&userDomain.Full_name,
		&userDomain.Email,
		&userDomain.Phone_number,
		&userDomain.Password,
		&userDomain.Time_add,
		&userDomain.Role,
		&userDomain.Is_verified)

	if err != nil {
		return core_domain.User{}, fmt.Errorf("error inserting user: %w", err)
	}

	return userDomain, nil
}

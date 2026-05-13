package features_admin_repository

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	"github.com/jackc/pgx/v5"
)

func (r *AdminRepository) GetUser(ctx context.Context, id int) (core_domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	query := `
		SELECT id, full_name, email, phone_number, password, role, time_add
		FROM trackerapp.users
		WHERE id = $1;`

	row := r.pool.QueryRow(ctx, query, id)

	var user core_domain.User
	var roleValue string
	err := row.Scan(
		&user.ID,
		&user.Full_name,
		&user.Email,
		&user.Phone_number,
		&user.Password,
		&roleValue,
		&user.Time_add,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_domain.User{}, core_errors.ErrorNotFoud
		}
		return core_domain.User{}, fmt.Errorf("get admin user scan: %w", err)
	}

	user.Role = core_domain.UserRole(roleValue)

	return user, nil
}

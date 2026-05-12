package features_users_repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	"github.com/jackc/pgx/v5"
)

func (r *UserRepository) ExtraditionUser(ctx context.Context, id int) (core_domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT * 
		FROM social.users 
        WHERE id=$1;`

	row := r.pool.QueryRow(ctx, query, id)

	var user core_domain.User
	err := row.Scan(
		&user.ID,
		&user.Full_name,
		&user.Email,
		&user.Phone_number,
		&user.Password,
		&user.Time_add)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_domain.User{}, core_errors.ErrorNotFoud
		}
		return core_domain.User{}, fmt.Errorf("ExtraditionUser scan: %w", err)
	}

	return user, nil
}

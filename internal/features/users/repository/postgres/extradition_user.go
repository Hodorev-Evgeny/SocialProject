package features_users_repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
)

func (r *UserRepository) ExtraditionUser(ctx context.Context, id int) (core_domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `SELECT * 
		FROM trackerapp.users 
        WHERE id=$1;`

	row, err := r.pool.QueryRow(ctx, query, id)
	if err != nil {
		return core_domain.User{}, fmt.Errorf("ExtraditionUser query: %w", err)
	}

	var user core_domain.User
	er := row.Scan(
		&user.ID,
		&user.Full_name,
		&user.Email,
		&user.Phone_number,
		&user.Password,
		&user.Time_add)
	if er != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core_domain.User{}, core_errors.ErrorNotFoud
		}
		return core_domain.User{}, fmt.Errorf("ExtraditionUser scan: %w", err)
	}

	return user, nil
}

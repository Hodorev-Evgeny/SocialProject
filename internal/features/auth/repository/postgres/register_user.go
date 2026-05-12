package features_auth_repository

import (
	"context"
	"time"
)

func (r *AuthRepository) InsertRegisteredUser(
	ctx context.Context,
	fullName, email string,
	phone *string,
	passwordHash string,
	timeAdd time.Time,
	role string,
) (userID int, err error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	query := `
	INSERT INTO trackerapp.users (full_name, email, phone_number, password, time_add, role, is_verified)
	VALUES ($1, $2, $3, $4, $5, $6, FALSE)
	RETURNING id;
	`

	row := r.pool.QueryRow(ctx, query,
		fullName,
		email,
		phone,
		passwordHash,
		timeAdd,
		role,
	)

	if err := row.Scan(&userID); err != nil {
		return 0, err
	}

	return userID, nil
}

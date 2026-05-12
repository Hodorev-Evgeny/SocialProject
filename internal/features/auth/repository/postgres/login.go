package features_auth_repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (r *AuthRepository) GetUserPasswordHashByEmail(ctx context.Context, email string) (userID int, passwordHash string, err error) {
	ctx, close := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer close()

	query := `
	SELECT id, password
	FROM trackerapp.users
	WHERE email = $1;
	`

	row := r.pool.QueryRow(ctx, query, email)
	if err := row.Scan(&userID, &passwordHash); err != nil {
		return 0, "", err
	}

	return userID, passwordHash, nil
}

func (r *AuthRepository) CreateVerification(ctx context.Context, id uuid.UUID, userID int, purpose string, codeHash string, expiresAt time.Time) error {
	ctx, close := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer close()

	query := `
	INSERT INTO trackerapp.auth_verifications (id, user_id, purpose, code_hash, expires_at)
	VALUES ($1, $2, $3, $4, $5);
	`

	_, err := r.pool.Exec(ctx, query, id, userID, purpose, codeHash, expiresAt)
	return err
}

package features_auth_repository

import (
	"context"
	"time"
)

func (r *AuthRepository) GetRefreshSessionByTokenHash(ctx context.Context, tokenHash string) (userID int, expiresAt time.Time, err error) {
	err = r.pool.QueryRow(ctx, `
		SELECT user_id, expires_at
		FROM trackerapp.refresh_sessions
		WHERE token_hash = $1
	`, tokenHash).Scan(&userID, &expiresAt)
	if err != nil {
		return 0, time.Time{}, err
	}
	return userID, expiresAt, nil
}

func (r *AuthRepository) DeleteRefreshSessionByTokenHash(ctx context.Context, tokenHash string) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM trackerapp.refresh_sessions WHERE token_hash = $1
	`, tokenHash)
	return err
}

package features_auth_repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *AuthRepository) InsertRevokedAccessToken(ctx context.Context, jti uuid.UUID, expiresAt time.Time) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO trackerapp.revoked_access_tokens (jti, expires_at)
		VALUES ($1, $2)
		ON CONFLICT (jti) DO NOTHING
	`, jti, expiresAt)
	return err
}

func (r *AuthRepository) IsAccessTokenRevoked(ctx context.Context, jti uuid.UUID) (bool, error) {
	var one int
	err := r.pool.QueryRow(ctx, `
		SELECT 1 FROM trackerapp.revoked_access_tokens WHERE jti = $1
	`, jti).Scan(&one)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

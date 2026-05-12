package features_auth_repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (r *AuthRepository) CreateRegisterVerification(ctx context.Context, id uuid.UUID, userID int, purpose string, codeHash string, expiresAt time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	query := `
	INSERT INTO trackerapp.register_verifications (id, user_id, purpose, code_hash, expires_at)
	VALUES ($1, $2, $3, $4, $5);
	`

	_, err := r.pool.Exec(ctx, query, id, userID, purpose, codeHash, expiresAt)
	return err
}

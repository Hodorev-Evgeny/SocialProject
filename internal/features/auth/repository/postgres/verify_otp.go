package features_auth_repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type VerifyOtpParams struct {
	VerificationID   uuid.UUID
	ExpectedPurpose  string
	PlainCode        string
	MarkUserVerified bool
	RefreshSessionID uuid.UUID
	RefreshTokenHash string
	RefreshExpiresAt time.Time
}

func (r *AuthRepository) VerifyOtpConsumeAndCreateRefreshSession(ctx context.Context, p VerifyOtpParams) (userID int, err error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var purpose string
	var codeHash string
	var expiresAt time.Time
	var consumedAt *time.Time

	row := tx.QueryRow(ctx, `
		SELECT user_id, purpose, code_hash, expires_at, consumed_at
		FROM trackerapp.auth_verifications
		WHERE id = $1
		FOR UPDATE;
	`, p.VerificationID)

	if err := row.Scan(&userID, &purpose, &codeHash, &expiresAt, &consumedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("%w: verification not found", core_errors.ErrorValidation)
		}
		return 0, fmt.Errorf("load verification: %w", err)
	}

	if consumedAt != nil {
		return 0, fmt.Errorf("%w: verification already used", core_errors.ErrorValidation)
	}

	if purpose != p.ExpectedPurpose {
		return 0, fmt.Errorf("%w: verification purpose mismatch", core_errors.ErrorValidation)
	}

	if time.Now().UTC().After(expiresAt.UTC()) {
		return 0, fmt.Errorf("%w: verification expired", core_errors.ErrorValidation)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(codeHash), []byte(p.PlainCode)); err != nil {
		return 0, fmt.Errorf("%w: invalid code", core_errors.ErrorValidation)
	}

	if _, err := tx.Exec(ctx, `
		UPDATE trackerapp.auth_verifications
		SET consumed_at = NOW()
		WHERE id = $1 AND consumed_at IS NULL;
	`, p.VerificationID); err != nil {
		return 0, fmt.Errorf("consume verification: %w", err)
	}

	if p.MarkUserVerified {
		if _, err := tx.Exec(ctx, `
			UPDATE trackerapp.users
			SET is_verified = TRUE
			WHERE id = $1;
		`, userID); err != nil {
			return 0, fmt.Errorf("mark user verified: %w", err)
		}
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO trackerapp.refresh_sessions (id, user_id, token_hash, expires_at)
		VALUES ($1, $2, $3, $4);
	`, p.RefreshSessionID, userID, p.RefreshTokenHash, p.RefreshExpiresAt); err != nil {
		return 0, fmt.Errorf("create refresh session: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit verify tx: %w", err)
	}

	return userID, nil
}

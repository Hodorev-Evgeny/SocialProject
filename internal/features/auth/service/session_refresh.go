package features_auth_service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func refreshOpaqueTokenHash(refreshRaw string) (string, error) {
	raw := strings.TrimSpace(refreshRaw)
	if raw == "" {
		return "", fmt.Errorf("%w: refresh_token is required", core_errors.ErrorValidation)
	}
	b, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return "", fmt.Errorf("%w: invalid refresh_token", core_errors.ErrorValidation)
	}
	if len(b) != 32 {
		return "", fmt.Errorf("%w: invalid refresh_token", core_errors.ErrorValidation)
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}

// Logout deletes the refresh session for refreshRaw. Missing or unknown refresh is idempotent (no error).
// If accessRawOptional is a valid access JWT for the same user as the refresh session, that access token
// is revoked immediately (subsequent /users/me with that JWT returns 401).
// If accessRawOptional is present and valid but for another user than the refresh session, returns ErrorValidation (422).
func (s *AuthService) Logout(ctx context.Context, accessRawOptional string, refreshRaw string) error {
	hash, err := refreshOpaqueTokenHash(refreshRaw)
	if err != nil {
		return err
	}

	var accessClaims *accessClaims
	if raw := strings.TrimSpace(accessRawOptional); raw != "" {
		if c, err := s.parseValidAccessClaims(raw); err == nil {
			accessClaims = c
		}
	}

	userID, _, err := s.auth.GetRefreshSessionByTokenHash(ctx, hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("lookup refresh session: %w", err)
	}
	if accessClaims != nil {
		accessUID, err := strconv.Atoi(accessClaims.Subject)
		if err != nil || accessUID <= 0 {
			return fmt.Errorf("%w: invalid access token subject", core_errors.ErrorValidation)
		}
		if accessUID != userID {
			return fmt.Errorf("%w: refresh token does not match authenticated user", core_errors.ErrorValidation)
		}
	}
	if err := s.auth.DeleteRefreshSessionByTokenHash(ctx, hash); err != nil {
		return fmt.Errorf("delete refresh session: %w", err)
	}

	if accessClaims != nil && accessClaims.ID != "" && accessClaims.ExpiresAt != nil {
		jti, err := uuid.Parse(accessClaims.ID)
		if err != nil {
			return nil
		}
		if err := s.auth.InsertRevokedAccessToken(ctx, jti, accessClaims.ExpiresAt.Time.UTC()); err != nil {
			return fmt.Errorf("revoke access token: %w", err)
		}
	}
	return nil
}

// Refresh issues a new access JWT; the refresh token string is unchanged (no rotation).
func (s *AuthService) Refresh(ctx context.Context, refreshRaw string) (IssuedTokenPair, error) {
	refreshRaw = strings.TrimSpace(refreshRaw)
	hash, err := refreshOpaqueTokenHash(refreshRaw)
	if err != nil {
		return IssuedTokenPair{}, err
	}
	userID, expiresAt, err := s.auth.GetRefreshSessionByTokenHash(ctx, hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return IssuedTokenPair{}, core_errors.ErrorUnauthorized
		}
		return IssuedTokenPair{}, fmt.Errorf("lookup refresh session: %w", err)
	}
	if time.Now().UTC().After(expiresAt.UTC()) {
		return IssuedTokenPair{}, core_errors.ErrorUnauthorized
	}

	u, err := s.users.ExtraditionUser(ctx, userID)
	if err != nil {
		return IssuedTokenPair{}, fmt.Errorf("load user: %w", err)
	}

	access, expiresIn, err := s.signAccessToken(u.ID, u.Email, u.Role)
	if err != nil {
		return IssuedTokenPair{}, err
	}

	return IssuedTokenPair{
		AccessToken:  access,
		RefreshToken: refreshRaw,
		ExpiresInSec: expiresIn,
		UserID:       u.ID,
		FullName:     u.Full_name,
		Email:        u.Email,
		Role:         u.Role,
	}, nil
}

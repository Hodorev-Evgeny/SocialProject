package features_auth_service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	core_auth "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/auth"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type accessClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func (s *AuthService) signAccessToken(userID int, email, role string) (token string, expiresInSec int64, err error) {
	now := time.Now().UTC()
	exp := now.Add(s.jwt.AccessTTL)
	jti := uuid.New().String()

	claims := accessClaims{
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(userID),
			ID:        jti,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString([]byte(s.jwt.Secret))
	if err != nil {
		return "", 0, fmt.Errorf("sign access token: %w", err)
	}

	return signed, int64(s.jwt.AccessTTL.Seconds()), nil
}

func (s *AuthService) parseValidAccessClaims(raw string) (*accessClaims, error) {
	var claims accessClaims
	token, err := jwt.ParseWithClaims(raw, &claims, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(s.jwt.Secret), nil
	})
	if err != nil || token == nil || !token.Valid || claims.Subject == "" {
		return nil, core_errors.ErrorUnauthorized
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil || userID <= 0 {
		return nil, core_errors.ErrorUnauthorized
	}

	return &claims, nil
}

// VerifyAccessToken валидирует access JWT и возвращает principal для мидлвари и хендлеров.
func (s *AuthService) VerifyAccessToken(ctx context.Context, raw string) (core_auth.Principal, error) {
	claims, err := s.parseValidAccessClaims(raw)
	if err != nil {
		return core_auth.Principal{}, err
	}

	if claims.ID != "" {
		jti, err := uuid.Parse(claims.ID)
		if err != nil {
			return core_auth.Principal{}, core_errors.ErrorUnauthorized
		}
		revoked, err := s.auth.IsAccessTokenRevoked(ctx, jti)
		if err != nil {
			return core_auth.Principal{}, fmt.Errorf("check access revocation: %w", err)
		}
		if revoked {
			return core_auth.Principal{}, core_errors.ErrorUnauthorized
		}
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil || userID <= 0 {
		return core_auth.Principal{}, core_errors.ErrorUnauthorized
	}

	return core_auth.Principal{
		UserID: userID,
		Email:  claims.Email,
		Role:   claims.Role,
	}, nil
}

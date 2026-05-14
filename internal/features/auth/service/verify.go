package features_auth_service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	features_auth_repository "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/auth/repository/postgres"
	"github.com/google/uuid"
)

type userReader interface {
	ExtraditionUser(ctx context.Context, id int) (core_domain.User, error)
}

type IssuedTokenPair struct {
	AccessToken      string
	RefreshToken     string
	ExpiresInSec     int64
	UserID           int
	FullName         string
	Email            string
	Role             string
}

func (s *AuthService) VerifyLogin(ctx context.Context, verificationID string, code string) (IssuedTokenPair, error) {
	return s.verifyOtpAndIssueTokens(ctx, verificationID, code, "login", false)
}

func (s *AuthService) VerifyRegister(ctx context.Context, verificationID string, code string) (IssuedTokenPair, error) {
	return s.verifyOtpAndIssueTokens(ctx, verificationID, code, "registration", true)
}

func (s *AuthService) verifyOtpAndIssueTokens(
	ctx context.Context,
	verificationID string,
	code string,
	expectedPurpose string,
	markUserVerified bool,
) (IssuedTokenPair, error) {
	vid, err := parseVerificationID(verificationID)
	if err != nil {
		return IssuedTokenPair{}, err
	}

	code = strings.TrimSpace(code)
	if c := utf8.RuneCountInString(code); c < 4 || c > 10 {
		return IssuedTokenPair{}, fmt.Errorf("%w: invalid code length", core_errors.ErrorValidation)
	}

	refreshRaw, refreshHash, err := newRefreshOpaqueToken()
	if err != nil {
		return IssuedTokenPair{}, fmt.Errorf("refresh token: %w", err)
	}

	refreshID := uuid.New()
	refreshExpires := time.Now().UTC().Add(s.jwt.RefreshTTL)

	userID, err := s.auth.VerifyOtpConsumeAndCreateRefreshSession(ctx, features_auth_repository.VerifyOtpParams{
		VerificationID:   vid,
		ExpectedPurpose:  expectedPurpose,
		PlainCode:        code,
		MarkUserVerified: markUserVerified,
		RefreshSessionID: refreshID,
		RefreshTokenHash: refreshHash,
		RefreshExpiresAt: refreshExpires,
	})
	if err != nil {
		return IssuedTokenPair{}, err
	}

	u, err := s.users.ExtraditionUser(ctx, userID)
	if err != nil {
		return IssuedTokenPair{}, fmt.Errorf("load user after verify: %w", err)
	}

	access, expiresIn, err := s.signAccessToken(u.ID, u.Email, u.Role)
	if err != nil {
		return IssuedTokenPair{}, err
	}

	return IssuedTokenPair{
		AccessToken:      access,
		RefreshToken:     refreshRaw,
		ExpiresInSec:     expiresIn,
		UserID:           u.ID,
		FullName:         u.Full_name,
		Email:            u.Email,
		Role:             u.Role,
	}, nil
}

func parseVerificationID(raw string) (uuid.UUID, error) {
	raw = strings.TrimSpace(raw)
	const prefix = "ver_"
	if !strings.HasPrefix(raw, prefix) {
		return uuid.Nil, fmt.Errorf("%w: invalid verification_id", core_errors.ErrorValidation)
	}
	id, err := uuid.Parse(strings.TrimPrefix(raw, prefix))
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: invalid verification_id", core_errors.ErrorValidation)
	}
	return id, nil
}

func newRefreshOpaqueToken() (raw string, sha256Hex string, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}
	sum := sha256.Sum256(b)
	return base64.RawURLEncoding.EncodeToString(b), hex.EncodeToString(sum[:]), nil
}

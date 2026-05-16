package features_auth_service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	features_auth_repository "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/auth/repository/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	auth  Repository
	users userReader
	jwt   JWTConfig
}

type Repository interface {
	GetUserPasswordHashByEmail(ctx context.Context, email string) (userID int, passwordHash string, err error)
	CreateVerification(ctx context.Context, id uuid.UUID, userID int, purpose string, codeHash string, expiresAt time.Time) error
	CreateRegisterVerification(ctx context.Context, id uuid.UUID, userID int, purpose string, codeHash string, expiresAt time.Time) error
	InsertRegisteredUser(ctx context.Context, fullName, email string, phone *string, passwordHash string, timeAdd time.Time, role string) (userID int, err error)
	VerifyOtpConsumeAndCreateRefreshSession(ctx context.Context, p features_auth_repository.VerifyOtpParams) (userID int, err error)
	GetRefreshSessionByTokenHash(ctx context.Context, tokenHash string) (userID int, expiresAt time.Time, err error)
	DeleteRefreshSessionByTokenHash(ctx context.Context, tokenHash string) error
	InsertRevokedAccessToken(ctx context.Context, jti uuid.UUID, expiresAt time.Time) error
	IsAccessTokenRevoked(ctx context.Context, jti uuid.UUID) (bool, error)
}

func NewAuthService(auth Repository, users userReader, jwt JWTConfig) *AuthService {
	return &AuthService{auth: auth, users: users, jwt: jwt}
}

func (s *AuthService) Register(ctx context.Context, req core_domain.User) (verificationID string, otpCode string, err error) {
	req.Email = strings.TrimSpace(req.Email)
	req.Full_name = strings.TrimSpace(req.Full_name)

	if err := req.Validate(); err != nil {
		return "", "", fmt.Errorf("%w: %w", core_errors.ErrorValidation, err)
	}

	if req.Role != "passenger" && req.Role != "driver" {
		return "", "", fmt.Errorf("%w: invalid role", core_errors.ErrorValidation)
	}

	if ln := len(req.Password); ln < 8 || ln > 32 {
		return "", "", fmt.Errorf("%w: invalid password length", core_errors.ErrorValidation)
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", fmt.Errorf("hash password: %w", err)
	}

	userID, err := s.auth.InsertRegisteredUser(ctx, req.Full_name, req.Email, req.Phone_number, string(hashBytes), req.Time_add, req.Role)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return "", "", fmt.Errorf("%w: email already registered", core_errors.ErrorBadRequest)
		}
		return "", "", fmt.Errorf("create user: %w", err)
	}

	code, err := generate6DigitCode()
	if err != nil {
		return "", "", fmt.Errorf("generate otp code: %w", err)
	}

	codeHashBytes, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return "", "", fmt.Errorf("hash otp code: %w", err)
	}

	id := uuid.New()
	expiresAt := time.Now().UTC().Add(15 * time.Minute)
	if err := s.auth.CreateRegisterVerification(ctx, id, userID, "registration", string(codeHashBytes), expiresAt); err != nil {
		return "", "", fmt.Errorf("create verification: %w", err)
	}

	return "ver_" + id.String(), code, nil
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (verificationID string, otpCode string, err error) {
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)

	if _, parseErr := mail.ParseAddress(email); parseErr != nil {
		return "", "", fmt.Errorf("%w: invalid email", core_errors.ErrorValidation)
	}
	if password == "" {
		return "", "", fmt.Errorf("%w: empty password", core_errors.ErrorValidation)
	}

	userID, passwordHash, err := s.auth.GetUserPasswordHashByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", core_errors.ErrorUnauthorized
		}
		return "", "", fmt.Errorf("get user by email: %w", err)
	}

	if compareErr := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); compareErr != nil {
		return "", "", core_errors.ErrorUnauthorized
	}

	code, err := generate6DigitCode()
	if err != nil {
		return "", "", fmt.Errorf("generate otp code: %w", err)
	}

	codeHashBytes, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return "", "", fmt.Errorf("hash otp code: %w", err)
	}

	id := uuid.New()
	expiresAt := time.Now().UTC().Add(15 * time.Minute)
	if err := s.auth.CreateVerification(ctx, id, userID, "login", string(codeHashBytes), expiresAt); err != nil {
		return "", "", fmt.Errorf("create verification: %w", err)
	}

	return "ver_" + id.String(), code, nil
}

func generate6DigitCode() (string, error) {
	// n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	// if err != nil {
	// 	return "", err
	// }
	// return fmt.Sprintf("%06d", n.Int64()), nil
	return "000000", nil
}

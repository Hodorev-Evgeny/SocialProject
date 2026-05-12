package features_auth_transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	features_auth_service "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/auth/service"
)

type OtpVerifyRequest struct {
	VerificationID string `json:"verification_id"`
	Code             string `json:"code"`
}

func (r *OtpVerifyRequest) Validate() error {
	if strings.TrimSpace(r.VerificationID) == "" {
		return fmt.Errorf("%w: verification_id is required", core_errors.ErrorValidation)
	}
	code := strings.TrimSpace(r.Code)
	if code == "" {
		return fmt.Errorf("%w: code is required", core_errors.ErrorValidation)
	}
	if c := utf8.RuneCountInString(code); c < 4 || c > 10 {
		return fmt.Errorf("%w: invalid code length", core_errors.ErrorValidation)
	}
	return nil
}

type tokenUserResponse struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type tokenPairResponse struct {
	AccessToken  string            `json:"access_token"`
	RefreshToken string            `json:"refresh_token"`
	TokenType    string            `json:"token_type"`
	ExpiresIn    int64             `json:"expires_in,omitempty"`
	User         tokenUserResponse `json:"user"`
}

func issuedToHTTP(t features_auth_service.IssuedTokenPair) tokenPairResponse {
	return tokenPairResponse{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		TokenType:    "bearer",
		ExpiresIn:    t.ExpiresInSec,
		User: tokenUserResponse{
			ID:       t.UserID,
			FullName: t.FullName,
			Email:    t.Email,
			Role:     t.Role,
		},
	}
}

func (h *AuthHTTPHandler) LoginVerify(w http.ResponseWriter, req *http.Request) {
	h.verifyOtp(w, req, h.authService.VerifyLogin)
}

func (h *AuthHTTPHandler) RegisterVerify(w http.ResponseWriter, req *http.Request) {
	h.verifyOtp(w, req, h.authService.VerifyRegister)
}

type verifyFn func(ctx context.Context, verificationID, code string) (features_auth_service.IssuedTokenPair, error)

func (h *AuthHTTPHandler) verifyOtp(w http.ResponseWriter, req *http.Request, fn verifyFn) {
	ctx := req.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := response.NewHandlerResponse(log, w)

	var body OtpVerifyRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		responseHandler.ErrorResponse(
			fmt.Errorf("%w: %w", core_errors.ErrorValidation, err),
			"error decoding request body",
		)
		return
	}

	body.VerificationID = strings.TrimSpace(body.VerificationID)
	body.Code = strings.TrimSpace(body.Code)

	if err := body.Validate(); err != nil {
		responseHandler.ErrorResponse(err, "validation error")
		return
	}

	tokens, err := fn(ctx, body.VerificationID, body.Code)
	if err != nil {
		responseHandler.ErrorResponse(err, "verify otp")
		return
	}

	responseHandler.JSONResponseHandler(http.StatusOK, issuedToHTTP(tokens))
}

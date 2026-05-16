package features_auth_transport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	"go.uber.org/zap"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	if strings.TrimSpace(r.Email) == "" {
		return fmt.Errorf("%w: email is required", core_errors.ErrorValidation)
	}
	if strings.TrimSpace(r.Password) == "" {
		return fmt.Errorf("%w: password is required", core_errors.ErrorValidation)
	}
	return nil
}

type AuthChallengeResponse struct {
	VerificationID string `json:"verification_id"`
	Message        string `json:"message"`
}

func (h *AuthHTTPHandler) Login(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := response.NewHandlerResponse(log, w)

	var body LoginRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		responseHandler.ErrorResponse(
			fmt.Errorf("%w: %w", core_errors.ErrorValidation, err),
			"error decoding request body",
		)
		return
	}
	if err := body.Validate(); err != nil {
		responseHandler.ErrorResponse(err, "validation error")
		return
	}

	verificationID, otpCode, err := h.authService.Login(ctx, body.Email, body.Password)
	if err != nil {
		responseHandler.ErrorResponse(err, "error login")
		return
	}

	log.Info(
		"login otp generated",
		zap.String("email", body.Email),
		zap.String("verification_id", verificationID),
		zap.String("code", otpCode),
	)

	if err := h.otpMailer.SendLoginOTP(ctx, body.Email, otpCode); err != nil {
		log.Error("login otp resend send failed", zap.Error(err), zap.String("email", body.Email))
	}

	responseHandler.JSONResponseHandler(http.StatusOK, AuthChallengeResponse{
		VerificationID: verificationID,
		Message:        "otp_sent",
	})
}

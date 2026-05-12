package features_auth_transport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	"go.uber.org/zap"
)

type RegisterRequest struct {
	FullName string  `json:"full_name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Role     string  `json:"role"`
	Phone    *string `json:"phone"`
}

func (r *RegisterRequest) Validate() error {
	fullName := strings.TrimSpace(r.FullName)
	if fullName == "" {
		return fmt.Errorf("%w: full_name is required", core_errors.ErrorValidation)
	}
	if rn := utf8.RuneCountInString(fullName); rn < 3 || rn > 32 {
		return fmt.Errorf("%w: full_name length must be between 3 and 32", core_errors.ErrorValidation)
	}

	if strings.TrimSpace(r.Email) == "" {
		return fmt.Errorf("%w: email is required", core_errors.ErrorValidation)
	}

	if strings.TrimSpace(r.Password) == "" {
		return fmt.Errorf("%w: password is required", core_errors.ErrorValidation)
	}
	if ln := utf8.RuneCountInString(r.Password); ln < 8 || ln > 32 {
		return fmt.Errorf("%w: password length must be between 8 and 32", core_errors.ErrorValidation)
	}

	if strings.TrimSpace(r.Role) == "" {
		return fmt.Errorf("%w: role is required", core_errors.ErrorValidation)
	}
	if r.Role != "passenger" && r.Role != "driver" {
		return fmt.Errorf("%w: invalid role", core_errors.ErrorValidation)
	}

	return nil
}

func (h *AuthHTTPHandler) Register(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := response.NewHandlerResponse(log, w)

	var body RegisterRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		responseHandler.ErrorResponse(
			fmt.Errorf("%w: %w", core_errors.ErrorValidation, err),
			"error decoding request body",
		)
		return
	}

	body.FullName = strings.TrimSpace(body.FullName)
	body.Email = strings.TrimSpace(body.Email)
	body.Role = strings.TrimSpace(body.Role)

	if body.Phone != nil {
		tr := strings.TrimSpace(*body.Phone)
		if tr == "" {
			body.Phone = nil
		} else {
			body.Phone = &tr
		}
	}

	if err := body.Validate(); err != nil {
		responseHandler.ErrorResponse(err, "validation error")
		return
	}

	domainUser := core_domain.CreateUnincelizedUser(
		body.FullName,
		body.Email,
		body.Phone,
		body.Password,
	)
	domainUser.Role = body.Role

	verificationID, otpCode, err := h.authService.Register(ctx, domainUser)
	if err != nil {
		responseHandler.ErrorResponse(err, "error register")
		return
	}

	log.Info(
		"register otp generated (email send stubbed)",
		zap.String("email", body.Email),
		zap.String("verification_id", verificationID),
		zap.String("code", otpCode),
	)

	responseHandler.JSONResponseHandler(http.StatusCreated, AuthChallengeResponse{
		VerificationID: verificationID,
		Message:        "otp_sent",
	})
}

package features_auth_transport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	core_middleware "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/middleware"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
)

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (r *refreshTokenRequest) Validate() error {
	if strings.TrimSpace(r.RefreshToken) == "" {
		return fmt.Errorf("%w: refresh_token is required", core_errors.ErrorValidation)
	}
	return nil
}

func (h *AuthHTTPHandler) Logout(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := response.NewHandlerResponse(log, w)

	var body refreshTokenRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		responseHandler.ErrorResponse(
			fmt.Errorf("%w: %w", core_errors.ErrorValidation, err),
			"error decoding request body",
		)
		return
	}
	body.RefreshToken = strings.TrimSpace(body.RefreshToken)
	if err := body.Validate(); err != nil {
		responseHandler.ErrorResponse(err, "validation error")
		return
	}

	var accessRaw string
	if raw, ok := core_middleware.RawBearerAccessToken(req); ok {
		accessRaw = raw
	}

	if err := h.authService.Logout(ctx, accessRaw, body.RefreshToken); err != nil {
		responseHandler.ErrorResponse(err, "logout")
		return
	}
	responseHandler.NoContent()
}

func (h *AuthHTTPHandler) Refresh(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := response.NewHandlerResponse(log, w)

	var body refreshTokenRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		responseHandler.ErrorResponse(
			fmt.Errorf("%w: %w", core_errors.ErrorValidation, err),
			"error decoding request body",
		)
		return
	}
	body.RefreshToken = strings.TrimSpace(body.RefreshToken)
	if err := body.Validate(); err != nil {
		responseHandler.ErrorResponse(err, "validation error")
		return
	}

	tokens, err := h.authService.Refresh(ctx, body.RefreshToken)
	if err != nil {
		responseHandler.ErrorResponse(err, "refresh token")
		return
	}

	responseHandler.JSONResponseHandler(http.StatusOK, issuedToHTTP(tokens))
}

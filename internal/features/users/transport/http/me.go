package features_users_transport

import (
	"net/http"

	core_auth "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/auth"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
)

func (h *UserHTTPHandler) CurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := response.NewHandlerResponse(log, w)

	p, ok := core_auth.PrincipalFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(core_errors.ErrorUnauthorized, "unauthorized")
		return
	}

	u, err := h.userService.ExtraditionUser(ctx, p.UserID)
	if err != nil {
		responseHandler.ErrorResponse(err, "error loading current user")
		return
	}

	responseHandler.JSONResponseHandler(http.StatusOK, DomainToCurrentUser(u))
}

package features_admin_transport

import (
	"net/http"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_utils "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

func (h *AdminHTTPHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := response.NewHandlerResponse(log, w)

	log.Info("Start processing get admin user id")
	id, err := core_http_utils.GetValuePathInt(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "id is required")
		return
	}

	log.Info("Start processing get user by admin")
	user, err := h.adminService.GetUser(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse(err, "error getting user by admin")
		return
	}
	log.Info("End processing get user by admin")

	responseHandler.JSONResponseHandler(http.StatusOK, DomainFromResponse(user))
}

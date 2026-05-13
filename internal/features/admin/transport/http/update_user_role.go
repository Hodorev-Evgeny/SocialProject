package features_admin_transport

import (
	"net/http"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_utils "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

type UpdateUserRoleRequest struct {
	Role string `json:"role"`
}

func (h *AdminHTTPHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := response.NewHandlerResponse(log, w)

	log.Info("Start processing get admin user id")
	id, err := core_http_utils.GetValuePathInt(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "id is required")
		return
	}

	log.Info("Start processing update user role request")
	var req UpdateUserRoleRequest
	if err := core_http_utils.DecodeJSON(&req, r); err != nil {
		responseHandler.ErrorResponse(err, "error parsing request")
		return
	}

	log.Info("Start processing update user role by admin")
	user, err := h.adminService.UpdateUserRole(ctx, id, core_domain.UserRole(req.Role))
	if err != nil {
		responseHandler.ErrorResponse(err, "error updating user role by admin")
		return
	}
	log.Info("End processing update user role by admin")

	responseHandler.JSONResponseHandler(http.StatusOK, DomainFromResponse(user))
}

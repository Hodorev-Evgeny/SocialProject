package features_users_transport

import (
	"net/http"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_query_parm "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

func (h *UserHTTPHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := response.NewHandlerResponse(log, w)

	log.Info("Start processing get value path")
	id, err := core_http_query_parm.GetValuePathInt(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "id is required")
		return
	}

	log.Info("End start processing delete user")
	if err := h.userService.DeleteUser(ctx, id); err != nil {
		responseHandler.ErrorResponse(err, err.Error())
		return
	}
	log.Info("End processing delete user")

	responseHandler.JSONResponseHandler(http.StatusNoContent, nil)
}

package features_users_transport

import (
	"net/http"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_query_parm "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

func (h *UserHTTPHandler) ExtraditionUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := response.NewHandlerResponse(log, w)

	log.Info("End processing get value path")
	id, err := core_http_query_parm.GetValuePathInt(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "id is required")
		return
	}

	log.Info("Start processing get user")
	domainUser, err := h.userService.ExtraditionUser(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse(err, "error try extradition user")
		return
	}
	log.Info("End processing get user")

	req := DomainFromResponse(domainUser)

	responseHandler.JSONResponseHandler(http.StatusOK, req)
}

package feature_transport_passenger

import (
	"net/http"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_query_parm "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

func (t *PassengerTransport) GetDriveByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	passengerId, err := core_http_query_parm.GetValuePathInt(r, "id")
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error getting passenger id")
		return
	}

	driveDomain, err := t.service.GetDrive(ctx, passengerId)
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error getting drive domain")
		return
	}

	resp := DomainDriveToResponse(driveDomain)

	ResponseHandler.JSONResponseHandler(http.StatusOK, resp)
}

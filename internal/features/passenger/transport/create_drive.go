package feature_transport_passenger

import (
	"net/http"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_query_parm "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

func (t *PassengerTransport) CreateDrive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	passengerID, err := core_http_query_parm.GetValuePathInt(r, "id")
	if err != nil {
		ResponseHandler.ErrorResponse(err, "id is required")
		return
	}

	domainPassenger, err := t.service.CreateDrive(ctx, passengerID)
	if err != nil {
		ResponseHandler.ErrorResponse(err, "failed to create drive")
		return
	}

	resp := DomainDriveToResponse(domainPassenger)
	ResponseHandler.JSONResponseHandler(http.StatusCreated, resp)
}

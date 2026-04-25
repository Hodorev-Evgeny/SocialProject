package features_users_transport

import (
	"fmt"
	"net/http"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_query_parm "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

type GetUsersResponse []UserDTOResponse

func (h *UserHTTPHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	limit, offset, err := GetLimitAnsOffset(r)
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error parsing limit/offset")
		return
	}

	usersDomain, err := h.userService.GetUsers(ctx, limit, offset)
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error getting users")
		return
	}

	userRsponse := DomainFromDTOesponse(usersDomain)

	ResponseHandler.JSONResponseHandler(http.StatusOK, userRsponse)
}

func GetLimitAnsOffset(r *http.Request) (*int, *int, error) {
	limit, err := core_http_query_parm.GetIntQueryParm(r, "limit")

	if err != nil {
		return nil, nil, fmt.Errorf("error parsing limit: %w", core_errors.ErrorValidation)
	}

	offset, err := core_http_query_parm.GetIntQueryParm(r, "offset")
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing offset: %w", core_errors.ErrorValidation)
	}

	return limit, offset, nil
}

func DomainFromDTOesponse(rows []core_domain.User) GetUsersResponse {
	users := make([]UserDTOResponse, 0)
	for _, user := range rows {
		userResponse := DomainFromResponse(user)
		users = append(users, userResponse)
	}
	return users
}

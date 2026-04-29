package features_users_transport

import (
	"errors"
	"net/http"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_types "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/types"
	core_http_utils "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
	"go.uber.org/zap"
)

type RequestPatchUser struct {
	FullName core_http_types.Nullable[string] `json:"full_name"`
	Email    core_http_types.Nullable[string] `json:"email"`
	Phone    core_http_types.Nullable[string] `json:"phone"`
}

func (u *RequestPatchUser) Validate() error {
	if u.FullName.Set {
		if u.FullName.Value == nil {
			return errors.New("full name must not be empty")
		}
	}
	if u.Email.Set {
		if u.Email.Value == nil {
			return errors.New("email must not be empty")
		}
	}

	return nil
}

func CreateUserPatch(user RequestPatchUser) core_domain.UserPatch {
	return core_domain.NewUserPatch(user.FullName.ToDomain(),
		user.Email.ToDomain(),
		user.Phone.ToDomain(),
	)
}

func (h *UserHTTPHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	log.Info("Start processing decoding body")
	var data RequestPatchUser
	if err := core_http_utils.DecodeJSON(&data, r); err != nil {
		ResponseHandler.ErrorResponse(err, "error parsing request")

		return
	}

	log.Info("Start processing patch user")
	userId, err := core_http_utils.GetValuePathInt(r, "id")
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error parsing id")
		return
	}

	log.Debug("PatchUser",
		zap.Any("data", data),
	)

	userPatch := CreateUserPatch(data)
	log.Info("End processing patch user")

	log.Info("Start processing response body")
	userDomain, err := h.userService.PatchUser(ctx, userId, userPatch)
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error patching user")
		return
	}

	ResponseHandler.JSONResponseHandler(http.StatusOK, userDomain)
}

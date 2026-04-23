package features_users_transport

import (
	"encoding/json"
	"net/http"
	"time"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	FullName string  `json:"full_name" validate:"required,min=3,max=32"`
	Password string  `json:"password" validate:"required,min=8,max=32"`
	Email    string  `json:"email" validate:"required,email"`
	Phone    *string `json:"phone" validate:"required,min=11,max=11, startswith=+"`
}

type CreateUserResponse struct {
	ID       int       `json:"id"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email"`
	Phone    *string   `json:"phone"`
	TimeAdd  time.Time `json:"time_add"`
}

func (h *UserHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	RsponceHandler := response.NewHandlerResponse(log, w)

	log.Info("CreateUser called")

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RsponceHandler.ErrorResponse(err, "error decoding request body")
		return
	}

	userDomain, err := h.userService.CreateUser(ctx, DTOFromDomain(req))
	if err != nil {
		RsponceHandler.ErrorResponse(err, "error creating user")
		return
	}

	respons := DomainFromResponse(userDomain)
	if err := json.NewEncoder(w).Encode(respons); err != nil {
		RsponceHandler.ErrorResponse(err, "error writing response")
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func DTOFromDomain(userDto CreateUserRequest) core_domain.User {
	return core_domain.CreateUnincelizedUser(
		userDto.FullName,
		userDto.Email,
		userDto.Phone,
		userDto.Password,
	)
}

func DomainFromResponse(user core_domain.User) CreateUserResponse {
	return CreateUserResponse{
		ID:       user.ID,
		FullName: user.Full_name,
		Phone:    user.Phone_number,
		TimeAdd:  user.Time_add,
	}
}

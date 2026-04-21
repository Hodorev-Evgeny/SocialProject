package features_users_transport

import (
	"net/http"

	core_transport_server "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/server"
)

type UserHTTPHandler struct {
	userService userService
}

type userService interface{}

func NewUserHTTPHandler(
	userService userService,
) *UserHTTPHandler {
	return &UserHTTPHandler{
		userService: userService,
	}
}

func (h *UserHTTPHandler) Routers() []core_transport_server.Route {
	return []core_transport_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
	}
}

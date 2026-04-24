package features_users_transport

import (
	"context"
	"net/http"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_transport_server "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/server"
)

type UserHTTPHandler struct {
	userService userService
}

type userService interface {
	CreateUser(ctx context.Context, req core_domain.User) (core_domain.User, error)
	GetUser(ctx context.Context, limit *int, offset *int) ([]core_domain.User, error)
}

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
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: h.GetUsers,
		},
	}
}

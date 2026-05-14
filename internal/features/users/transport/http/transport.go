package features_users_transport

import (
	"context"
	"net/http"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_middleware "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/middleware"
	core_transport_server "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/server"
)

type UserHTTPHandler struct {
	userService userService
	meHandler   http.Handler
}

type userService interface {
	CreateUser(
		ctx context.Context,
		req core_domain.User,
	) (core_domain.User, error)

	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]core_domain.User, error)

	ExtraditionUser(
		ctx context.Context,
		id int,
	) (core_domain.User, error)

	DeleteUser(
		ctx context.Context,
		id int,
	) error

	PatchUser(
		ctx context.Context,
		id int,
		req core_domain.UserPatch,
	) (core_domain.User, error)

	GetPassengerForDriver(
		ctx context.Context,
		id int,
	) (core_domain.User, error)
}

func NewUserHTTPHandler(
	userService userService,
	verifier core_middleware.AccessTokenVerifier,
) *UserHTTPHandler {
	h := &UserHTTPHandler{
		userService: userService,
	}
	h.meHandler = core_middleware.BearerAuth(verifier, http.HandlerFunc(h.CurrentUser))
	return h
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
		{
			Method:  http.MethodGet,
			Path:    "/users/me",
			Handler: http.HandlerFunc(h.meHandler.ServeHTTP),
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/{id}",
			Handler: h.ExtraditionUser,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/users/{id}",
			Handler: h.DeleteUser,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/users/{id}",
			Handler: h.PatchUser,
		},
		{
			Method:  http.MethodGet,
			Path:    "/drivers/passengers/{id}",
			Handler: h.GetPassengerForDriver,
		},
	}
}

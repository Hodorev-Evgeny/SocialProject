package features_admin_transport

import (
	"context"
	"net/http"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_transport_server "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/server"
)

type AdminHTTPHandler struct {
	adminService adminService
}

type adminService interface {
	GetUser(
		ctx context.Context,
		id int,
	) (core_domain.User, error)

	UpdateUserRole(
		ctx context.Context,
		id int,
		role core_domain.UserRole,
	) (core_domain.User, error)
}

func NewAdminHTTPHandler(adminService adminService) *AdminHTTPHandler {
	return &AdminHTTPHandler{
		adminService: adminService,
	}
}

func (h *AdminHTTPHandler) Routers() []core_transport_server.Route {
	return []core_transport_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/admin/users/{id}",
			Handler: h.GetUser,
		},
		{
			Method:  http.MethodPut,
			Path:    "/admin/users/{id}/role",
			Handler: h.UpdateUserRole,
		},
	}
}

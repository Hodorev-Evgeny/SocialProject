package feature_transport_passenger

import (
	"context"
	"net/http"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_transport_server "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/server"
)

type PassengerService interface {
	GetDrive(
		ctx context.Context,
		passengerID int,
	) (core_domain.Drive, error)

	TakeStop(
		ctx context.Context,
		passengerID int,
	) error

	CreateDrive(
		ctx context.Context,
		passengerID int,
	) (core_domain.Drive, error)
}

type PassengerTransport struct {
	service PassengerService
}

func NewPassengerTransport(
	service PassengerService,
) *PassengerTransport {
	return &PassengerTransport{
		service: service,
	}
}

func (t *PassengerTransport) Route() []core_transport_server.Route {
	return []core_transport_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/drives/{id}",
			Handler: t.CreateDrive,
		},
		{
			Method:  http.MethodGet,
			Path:    "/drives/{id}",
			Handler: t.GetDriveByID,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/drives/{id}",
			Handler: t.TakeStop,
		},
	}
}

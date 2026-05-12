package feature_service_passenger

import (
	"context"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

type PassengerRepository interface {
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
		drive core_domain.Drive,
	) (core_domain.Drive, error)
}

type PassengerService struct {
	PassengerRepository PassengerRepository
}

func CreatePassengerService(
	repository PassengerRepository,
) *PassengerService {
	return &PassengerService{
		PassengerRepository: repository,
	}
}

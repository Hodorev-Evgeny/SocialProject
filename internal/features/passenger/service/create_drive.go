package feature_service_passenger

import (
	"context"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (s *PassengerService) CreateDrive(
	ctx context.Context,
	passengerID int,
) (core_domain.Drive, error) {
	drive := core_domain.CreateDrive(passengerID)

	corectDrive, err := s.PassengerRepository.CreateDrive(ctx, drive)
	if err != nil {
		return core_domain.Drive{}, err
	}

	return corectDrive, nil
}

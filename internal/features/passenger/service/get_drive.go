package feature_service_passenger

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
)

func (s *PassengerService) GetDrive(
	ctx context.Context,
	passengerID int,
) (core_domain.Drive, error) {
	if passengerID == core_domain.UnincelizedID {
		return core_domain.Drive{}, core_errors.ErrorUnauthorized
	}

	drive, err := s.PassengerRepository.GetDrive(ctx, passengerID)
	if err != nil {
		return core_domain.Drive{}, fmt.Errorf("error getting drive from service: %w", err)
	}

	return drive, nil
}

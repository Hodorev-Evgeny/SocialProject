package feature_service_passenger

import (
	"context"
	"fmt"
)

func (s *PassengerService) TakeStop(
	ctx context.Context,
	passengerID int,
) error {
	drive, err := s.GetDrive(ctx, passengerID)
	if err != nil {
		return fmt.Errorf("error getting drive from service: %w", err)
	}

	if err := drive.FinishDrive(); err != nil {
		return fmt.Errorf("error finishing drive: %w", err)
	}

	if err := s.PassengerRepository.TakeStop(ctx, passengerID); err != nil {
		return fmt.Errorf("error taking stop: %w", err)
	}

	return nil
}

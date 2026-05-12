package feature_user_service

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
)

func (s *UserService) GetPassengerForDriver(ctx context.Context, id int) (core_domain.User, error) {
	if id == core_domain.UnincelizedID {
		return core_domain.User{}, fmt.Errorf("not allowed: %w", core_errors.ErrorValidation)
	}

	userDomain, err := s.userRepository.ExtraditionUser(ctx, id)
	if err != nil {
		return core_domain.User{}, fmt.Errorf("GetPassengerForDriver: %w", err)
	}

	return userDomain, nil
}

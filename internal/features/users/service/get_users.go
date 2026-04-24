package feature_user_service

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
)

func (s *UserService) GetUser(ctx context.Context, limit *int, offset *int) ([]core_domain.User, error) {
	if limit != nil && *limit < 0 {
		return []core_domain.User{},
			fmt.Errorf("limit is not valid %w", core_errors.ErrorValidation)
	}
	if offset != nil && *offset < 0 {
		return []core_domain.User{},
			fmt.Errorf("offset in not valid: %w", core_errors.ErrorValidation)
	}

	usersDomains, err := s.userRepository.GetUser(ctx, limit, offset)
	if err != nil {
		return []core_domain.User{}, fmt.Errorf("error getting users: %w", err)
	}

	return usersDomains, nil
}

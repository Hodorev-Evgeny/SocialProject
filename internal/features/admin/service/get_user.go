package feature_admin_service

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
)

func (s *AdminService) GetUser(ctx context.Context, id int) (core_domain.User, error) {
	if id == core_domain.UnincelizedID {
		return core_domain.User{}, fmt.Errorf("user id is not allowed: %w", core_errors.ErrorValidation)
	}

	user, err := s.adminRepository.GetUser(ctx, id)
	if err != nil {
		return core_domain.User{}, fmt.Errorf("get user by admin: %w", err)
	}

	return user, nil
}

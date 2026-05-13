package feature_admin_service

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
)

func (s *AdminService) UpdateUserRole(
	ctx context.Context,
	id int,
	role core_domain.UserRole,
) (core_domain.User, error) {
	if id == core_domain.UnincelizedID {
		return core_domain.User{}, fmt.Errorf("user id is not allowed: %w", core_errors.ErrorValidation)
	}

	if !role.IsAssignableByAdmin() {
		return core_domain.User{}, fmt.Errorf("role must be passenger or driver: %w", core_errors.ErrorValidation)
	}

	user, err := s.adminRepository.UpdateUserRole(ctx, id, role)
	if err != nil {
		return core_domain.User{}, fmt.Errorf("update user role by admin: %w", err)
	}

	return user, nil
}

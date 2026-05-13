package feature_admin_service

import (
	"context"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

type AdminService struct {
	adminRepository AdminRepository
}

type AdminRepository interface {
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

func NewAdminService(adminRepository AdminRepository) *AdminService {
	return &AdminService{
		adminRepository: adminRepository,
	}
}

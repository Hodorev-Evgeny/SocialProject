package feature_user_service

import (
	"context"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

type UserService struct {
	userRepository UserRepository
}

type UserRepository interface {
	AddUser(ctx context.Context, user core_domain.User) (core_domain.User, error)
	GetUsers(ctx context.Context, limit *int, offset *int) ([]core_domain.User, error)
	ExtraditionUser(ctx context.Context, id int) (core_domain.User, error)
}

func NewUserService(
	userRepository UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

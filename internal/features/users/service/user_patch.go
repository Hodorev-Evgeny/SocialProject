package feature_user_service

import (
	"context"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (s *UserService) PatchUser(ctx context.Context,
	id int,
	req core_domain.UserPatch,
) (core_domain.User, error) {
	user, err := s.ExtraditionUser(ctx, id)
	if err != nil {
		return core_domain.User{}, err
	}

	if err := user.ApplyPatch(req); err != nil {
		return core_domain.User{}, err
	}

	userUpdate, err := s.userRepository.PatchUser(ctx, id, user)
	if err != nil {
		return core_domain.User{}, err
	}

	return userUpdate, nil
}

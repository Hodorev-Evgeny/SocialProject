package feature_user_service

import (
	"context"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	if id == core_domain.UnincelizedID {
		return fmt.Errorf("cannot delete unincelized user")
	}

	err := s.userRepository.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

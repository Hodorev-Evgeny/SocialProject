package features_users_repository

import core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"

type UserRepository struct {
	pool core_repository_pool.Pool
}

func NewUserRepository(pool core_repository_pool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

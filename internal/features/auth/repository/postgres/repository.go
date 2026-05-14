package features_auth_repository

import core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"

type AuthRepository struct {
	pool core_repository_pool.Pool
}

func NewAuthRepository(pool core_repository_pool.Pool) *AuthRepository {
	return &AuthRepository{pool: pool}
}

package features_admin_repository

import core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"

type AdminRepository struct {
	pool core_repository_pool.Pool
}

func NewAdminRepository(pool core_repository_pool.Pool) *AdminRepository {
	return &AdminRepository{
		pool: pool,
	}
}

package feature_repository_passenger

import (
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
)

type PassengerRepository struct {
	pool core_repository_pool.Pool
}

func CreatePassengerRepository(
	pool core_repository_pool.Pool,
) *PassengerRepository {
	return &PassengerRepository{
		pool: pool,
	}
}

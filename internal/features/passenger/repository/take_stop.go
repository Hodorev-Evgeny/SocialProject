package feature_repository_passenger

import (
	"context"
	"errors"
	"time"

	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
)

func (r *PassengerRepository) TakeStop(
	ctx context.Context,
	passengerID int,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	query := `
		UPDATE social.drive
		SET status = $1, time_from = $2
		WHERE id = $3;
`

	timeFrom := time.Now().UTC()
	_, err := r.pool.Exec(ctx, query, "Finished", timeFrom, passengerID)
	if err != nil {
		if errors.Is(err, core_repository_pool.ErrNoRows) {
			return core_errors.ErrorNotFoud
		}

		return err
	}

	return nil
}

package feature_repository_passenger

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_repository_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool"
)

func (r *PassengerRepository) GetDrive(
	ctx context.Context,
	passengerID int,
) (core_domain.Drive, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	query := `
		SELECT id, user_id, time_to, time_from, status
		FROM social.drive
		WHERE id = $1;`

	row := r.pool.QueryRow(ctx, query, passengerID)

	var drive core_domain.Drive
	err := row.Scan(
		&drive.ID,
		&drive.UserID,
		&drive.TimeTo,
		&drive.TimeFrom,
		&drive.Status,
	)
	if err != nil {
		if errors.Is(err, core_repository_pool.ErrNoRows) {
			return core_domain.Drive{}, core_errors.ErrorNotFoud
		}
		return core_domain.Drive{}, fmt.Errorf("error getting drive from service: %w", err)
	}

	return drive, nil
}

package feature_repository_passenger

import (
	"context"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (r *PassengerRepository) CreateDrive(
	ctx context.Context,
	drive core_domain.Drive,
) (core_domain.Drive, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	query := `
		INSERT INTO social.drive (user_id, time_to, time_from, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, time_to, time_from, status;`

	row := r.pool.QueryRow(ctx, query,
		drive.UserID,
		drive.TimeTo,
		drive.TimeFrom,
		drive.Status,
	)

	var coreDrive core_domain.Drive
	err := row.Scan(
		&coreDrive.ID,
		&coreDrive.UserID,
		&coreDrive.TimeTo,
		&coreDrive.TimeFrom,
		&coreDrive.Status,
	)
	if err != nil {
		return core_domain.Drive{}, err
	}

	return coreDrive, nil
}

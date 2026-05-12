package feature_transport_passenger

import (
	"time"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

type ResponseDrive struct {
	ID       int        `json:"id"`
	UserID   int        `json:"user_id"`
	TimeTo   time.Time  `json:"time_to"`
	TimeFrom *time.Time `json:"time_from"`
	Status   string     `json:"status"`
}

func DomainDriveToResponse(
	drive core_domain.Drive,
) ResponseDrive {
	return ResponseDrive{
		ID:       drive.ID,
		UserID:   drive.UserID,
		TimeTo:   drive.TimeTo,
		TimeFrom: drive.TimeFrom,
		Status:   drive.Status,
	}
}

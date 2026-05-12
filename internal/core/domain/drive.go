package core_domain

import "time"

type Drive struct {
	ID       int
	UserID   int
	TimeTo   time.Time
	TimeFrom *time.Time
	Status   string
}

func (d *Drive) FinishDrive() error {
	d.Status = "Finished"
	timeFrom := time.Now()
	d.TimeFrom = &timeFrom
	return nil
}

func CreateDrive(userID int) Drive {
	timeTo := time.Now()
	return Drive{
		ID:       UnincelizedID,
		UserID:   userID,
		TimeTo:   timeTo,
		TimeFrom: nil,
		Status:   "Active",
	}
}

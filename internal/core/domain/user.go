package core_domain

import "time"

// создания доменой сущности пользователя
type User struct {
	ID           int
	Full_name    string
	Phone_number *string
	Password     string
	time_add     time.Time
}

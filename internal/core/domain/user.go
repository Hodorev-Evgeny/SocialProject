package core_domain

import (
	"time"
)

var (
	UnincelizedID = -1
)

// создания доменой сущности пользователя
type User struct {
	ID           int
	Full_name    string
	Email        string
	Phone_number *string
	Password     string
	Time_add     time.Time
}

func CreateUser(
	id int,
	full_name string,
	email string,
	phone *string,
	password string,
) User {
	return User{
		ID:           id,
		Full_name:    full_name,
		Email:        email,
		Phone_number: phone,
		Password:     password,
		Time_add:     time.Now(),
	}
}

func CreateUnincelizedUser(
	full_name string,
	email string,
	phone *string,
	password string) User {
	return CreateUser(
		UnincelizedID,
		full_name,
		email,
		phone,
		password,
	)
}

func (u *User) Validate() error {
	// сделать волидацию
	return nil
}

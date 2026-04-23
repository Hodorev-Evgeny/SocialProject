package core_domain

import (
	"fmt"
	"net/mail"
	"regexp"
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
	if u.Full_name == "" {
		return fmt.Errorf("invalid user full_name")
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return fmt.Errorf("invalid user email")
	}

	if u.Password == "" {
		return fmt.Errorf("invalid user password for user")
	}

	if u.Phone_number != nil {
		regular := regexp.MustCompile(`^\+[0-9]+$`)
		if !regular.MatchString(*u.Phone_number) {
			return fmt.Errorf("invalid user phone number")
		}
	}

	return nil
}

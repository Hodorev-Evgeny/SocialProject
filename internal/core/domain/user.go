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
	Role         string
	Is_verified  bool
	Time_add     time.Time
	Description  *string
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

	if u.Role != "" && u.Role != "passenger" && u.Role != "driver" {
		return fmt.Errorf("invalid user role")
	}

	if u.Phone_number != nil {
		regular := regexp.MustCompile(`^\+[0-9]+$`)
		if !regular.MatchString(*u.Phone_number) {
			return fmt.Errorf("invalid user phone number")
		}
	}

	return nil
}

type UserPatch struct {
	Full_name Nullable[string]
	Email     Nullable[string]
	Phone     Nullable[string]
}

func NewUserPatch(full_name Nullable[string],
	email Nullable[string],
	phone Nullable[string],
) UserPatch {
	return UserPatch{
		Full_name: full_name,
		Email:     email,
		Phone:     phone,
	}
}

func (u *UserPatch) Validate() error {
	if u.Full_name.Set && u.Full_name.Value == nil {
		return fmt.Errorf("invalid user full_name")
	}

	if u.Email.Set && u.Email.Value == nil {
		return fmt.Errorf("invalid user email")
	}

	if u.Phone.Set && u.Phone.Value != nil {
		regular := regexp.MustCompile(`^\+[0-9]+$`)
		if !regular.MatchString(*u.Phone.Value) {
			return fmt.Errorf("invalid user phone number")
		}
	}

	return nil
}

func (u *User) ApplyPatch(patch UserPatch) error {
	if patch.Validate() != nil {
		return fmt.Errorf("invalid user patch")
	}

	tmp := *u
	if patch.Full_name.Set {
		tmp.Full_name = *patch.Full_name.Value
	}

	if patch.Email.Set {
		tmp.Email = *patch.Email.Value
	}

	if patch.Phone.Set {
		tmp.Phone_number = patch.Phone.Value
	}

	if tmp.Validate() != nil {
		return fmt.Errorf("invalid user patch")
	}

	*u = tmp
	return nil
}

package features_users_transport

import (
	"time"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

type UserDTOResponse struct {
	ID          int       `json:"id"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	Phone       *string   `json:"phone"`
	TimeAdd     time.Time `json:"time_add"`
	Description *string   `json:"description"`
}

func DTOFromDomain(userDto CreateUserRequest) core_domain.User {
	return core_domain.CreateUnincelizedUser(
		userDto.FullName,
		userDto.Email,
		userDto.Phone,
		userDto.Password,
	)
}

func DomainFromResponse(user core_domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:          user.ID,
		FullName:    user.Full_name,
		Email:       user.Email,
		Phone:       user.Phone_number,
		TimeAdd:     user.Time_add,
		Description: user.Description,
	}
}

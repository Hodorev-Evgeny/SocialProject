package features_admin_transport

import (
	"time"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

type UserDTOResponse struct {
	ID       int       `json:"id"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email"`
	Phone    *string   `json:"phone"`
	Role     string    `json:"role"`
	TimeAdd  time.Time `json:"time_add"`
}

func DomainFromResponse(user core_domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:       user.ID,
		FullName: user.Full_name,
		Email:    user.Email,
		Phone:    user.Phone_number,
		Role:     string(user.Role),
		TimeAdd:  user.Time_add,
	}
}

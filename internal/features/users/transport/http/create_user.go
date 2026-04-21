package features_users_transport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
)

type CreateUserRequest struct {
	FullName string  `json:"full_name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Phone    *string `json:"phone"`
}

type CreateUserResponse struct {
	ID       string    `json:"id"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email"`
	Phone    *string   `json:"phone"`
	TimeAdd  time.Time `json:"time_add"`
}

func (h *UserHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	log.Info("CreateUser called")

	var req CreateUserRequest
	if err := json.NewEncoder(w).Encode(req); err != nil {
		fmt.Println("пу пу пу")
	}

	w.WriteHeader(http.StatusOK)
}

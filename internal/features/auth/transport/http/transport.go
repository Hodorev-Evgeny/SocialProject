package features_auth_transport

import (
	"context"
	"net/http"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_transport_server "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/server"
	features_auth_service "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/auth/service"
)

type AuthHTTPHandler struct {
	authService authService
}

type authService interface {
	Login(
		ctx context.Context,
		email string,
		password string,
	) (verificationID string, otpCode string, err error)

	Register(
		ctx context.Context,
		req core_domain.User,
	) (verificationID string, otpCode string, err error)

	VerifyLogin(
		ctx context.Context,
		verificationID string,
		code string,
	) (features_auth_service.IssuedTokenPair, error)

	VerifyRegister(
		ctx context.Context,
		verificationID string,
		code string,
	) (features_auth_service.IssuedTokenPair, error)
}

func NewAuthHTTPHandler(
	authService authService,
) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		authService: authService,
	}
}

func (h *AuthHTTPHandler) Routers() []core_transport_server.Route {
	return []core_transport_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/auth/login/verify",
			Handler: h.LoginVerify,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/login",
			Handler: h.Login,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/register/verify",
			Handler: h.RegisterVerify,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/register",
			Handler: h.Register,
		},
	}
}

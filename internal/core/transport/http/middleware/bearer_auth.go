package core_middleware

import (
	"net/http"
	"strings"

	core_auth "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/auth"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
)

// проверяет сырой access token и возвращает principal.
// Реализуется сервисом авторизации (например *AuthService).
type AccessTokenVerifier interface {
	VerifyAccessToken(rawToken string) (core_auth.Principal, error)
}

//	проверяет заголовок Authorization: Bearer <token> и кладёт
//
// core_auth.Principal в контекст запроса. Там где нужно проверять авторизацию пускать через него
func BearerAuth(verifier AccessTokenVerifier, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, err := parseBearerToken(r.Header.Get("Authorization"))
		if err != nil {
			writeUnauthorized(w, r, "missing or invalid authorization header")
			return
		}

		p, err := verifier.VerifyAccessToken(raw)
		if err != nil {
			writeUnauthorized(w, r, "invalid or expired access token")
			return
		}

		next.ServeHTTP(w, r.WithContext(core_auth.ContextWithPrincipal(r.Context(), p)))
	})
}

func parseBearerToken(header string) (string, error) {
	header = strings.TrimSpace(header)
	if header == "" {
		return "", core_errors.ErrorUnauthorized
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", core_errors.ErrorUnauthorized
	}
	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", core_errors.ErrorUnauthorized
	}
	return token, nil
}

func writeUnauthorized(w http.ResponseWriter, r *http.Request, msg string) {
	log := core_logger.FromContext(r.Context())
	response.NewHandlerResponse(log, w).ErrorResponse(core_errors.ErrorUnauthorized, msg)
}

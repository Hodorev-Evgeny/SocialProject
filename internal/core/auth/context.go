package core_auth

import "context"

// Principal — данные субъекта из access JWT (для контекста запроса).
type Principal struct {
	UserID int
	Email  string
	Role   string
}

type ctxKey int

const principalKey ctxKey = 1

func ContextWithPrincipal(ctx context.Context, p Principal) context.Context {
	return context.WithValue(ctx, principalKey, p)
}

func PrincipalFromContext(ctx context.Context) (Principal, bool) {
	p, ok := ctx.Value(principalKey).(Principal)
	return p, ok
}

package middleware

import (
	"context"

	"github.com/base-go/backend/pkg/response"
)

// SetUserContext stores the authenticated user on the request context.
// It is the single write site used by JWTAuthMiddleware.
func SetUserContext(ctx context.Context, user response.UserContext) context.Context {
	return context.WithValue(ctx, ContextUser, user)
}

// GetUserContext retrieves the authenticated user stored by
// JWTAuthMiddleware. Handlers must use this instead of reading the
// context with an ad-hoc key.
func GetUserContext(ctx context.Context) (response.UserContext, bool) {
	user, ok := ctx.Value(ContextUser).(response.UserContext)
	return user, ok
}

package middleware

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/base-go/backend/pkg/response"
)

func TestUserContextRoundTrip(t *testing.T) {
	userCtx := response.UserContext{
		UserID: "d0793289-7d71-48eb-826b-d5ea9648c1c6",
		Email:  "test@pollyglot.dev",
		Role:   "Member",
	}

	ctx := SetUserContext(context.Background(), userCtx)

	got, ok := GetUserContext(ctx)
	require.True(t, ok, "user context set by the middleware must be readable by handlers")
	assert.Equal(t, userCtx, got)
}

func TestGetUserContextMissing(t *testing.T) {
	_, ok := GetUserContext(context.Background())
	assert.False(t, ok)
}

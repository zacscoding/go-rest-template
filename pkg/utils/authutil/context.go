package authutil

import (
	"context"

	"github.com/gin-gonic/gin"
)

const (
	IdentityKey = "user.email"
)

type userContextKey string

const userKey = userContextKey("user")

type Principal interface {
	GetID() uint
	ToModel() interface{}
}

func CurrentUser(ctx context.Context) Principal {
	if ctx == nil {
		return nil
	}
	if p, ok := ctx.Value(userKey).(Principal); ok {
		return p
	}
	if gctx, ok := ctx.(*gin.Context); ok && gctx != nil {
		return currentUserFromGinContext(gctx)
	}
	return nil
}

func currentUserFromGinContext(ctx *gin.Context) Principal {
	v, ok := ctx.Get(IdentityKey)
	if !ok {
		return nil
	}
	if p, ok := v.(Principal); ok {
		return p
	}
	return nil
}

// WithUserContext creates a new context with the provided user attached.
func WithUserContext(ctx context.Context, p Principal) context.Context {
	return context.WithValue(ctx, userKey, p)
}

package auth

import (
	"context"
)

type ctxKey int

// key usada para armazenar e recuperar um Claims de um context.Context
const claimKey ctxKey = 1

// =============================================================================

func SetClaims(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, claimKey, claims)
}

func GetClaims(ctx context.Context) Claims {
	v, ok := ctx.Value(claimKey).(Claims)
	if !ok {
		return Claims{}
	}
	return v
}

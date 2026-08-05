package model

import "context"

type (
	userCtxKey  struct{}
	tokenCtxKey struct{}
)

func WithUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, userCtxKey{}, user)
}

func UserFromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(userCtxKey{}).(User)
	return user, ok
}

func WithToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenCtxKey{}, token)
}

func TokenFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(tokenCtxKey{}).(string)
	return token, ok
}

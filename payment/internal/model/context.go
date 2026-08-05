package model

import "context"

type userCtxKey struct{}

func WithUserContext(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, userCtxKey{}, user)
}

func UserFromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(userCtxKey{}).(User)
	return user, ok
}

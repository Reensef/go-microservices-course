// Package authinterceptor содержит общую реализацию gRPC unary-интерцептора,
// проверяющего Bearer-токен через IAM и кладущего пользователя в контекст.
// Тип пользователя у каждого сервиса свой, поэтому интерцептор параметризован
// IAM-клиентом и функцией добавления пользователя в контекст.
package authinterceptor

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const bearerPrefix = "Bearer "

var exemptMethods = map[string]struct{}{
	"/grpc.health.v1.Health/Check":                                   {},
	"/grpc.health.v1.Health/Watch":                                   {},
	"/grpc.reflection.v1.ServerReflection/ServerReflectionInfo":      {},
	"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo": {},
}

// IAMClient — минимальный контракт, нужный интерцептору для проверки сессии.
type IAMClient[T any] interface {
	Whoami(ctx context.Context, sessionUUID string) (T, error)
}

// New создаёт auth-интерцептор. withUserContext кладёт полученного пользователя в контекст,
// errInvalidSession — ошибка IAMClient, означающая невалидную/истёкшую сессию.
func New[T any](
	iamClient IAMClient[T],
	withUserContext func(ctx context.Context, user T) context.Context,
	errInvalidSession error,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if _, ok := exemptMethods[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing authorization metadata")
		}

		values := md.Get("authorization")
		if len(values) == 0 || !strings.HasPrefix(values[0], bearerPrefix) {
			return nil, status.Error(
				codes.Unauthenticated, "missing or malformed authorization metadata",
			)
		}

		token := strings.TrimPrefix(values[0], bearerPrefix)
		if token == "" {
			return nil, status.Error(
				codes.Unauthenticated, "missing or malformed authorization metadata",
			)
		}

		user, err := iamClient.Whoami(ctx, token)
		if err != nil {
			if errors.Is(err, errInvalidSession) {
				return nil, status.Error(codes.Unauthenticated, "invalid or expired session")
			}

			return nil, status.Error(codes.Internal, "failed to validate session")
		}

		return handler(withUserContext(ctx, user), req)
	}
}

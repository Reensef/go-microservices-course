// Package iamclient содержит общую реализацию gRPC-клиента к IAM-сервису,
// используемую inventory/order/payment. Тип пользователя у каждого сервиса свой,
// поэтому клиент параметризован функцией конвертации ответа IAM в локальную модель.
package iamclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	iamGrpc "github.com/Reensef/go-microservices-course/shared/pkg/proto/iam/v1"
)

// Client — gRPC-клиент AuthService IAM, конвертирующий ответ в локальный тип пользователя T.
type Client[T any] struct {
	service           iamGrpc.AuthServiceClient
	convert           func(*iamGrpc.User) T
	errInvalidSession error
}

// New создаёт клиент. convert конвертирует proto-пользователя в локальную модель сервиса,
// errInvalidSession — ошибка, которую нужно вернуть при истёкшей/неверной сессии.
func New[T any](service iamGrpc.AuthServiceClient, convert func(*iamGrpc.User) T, errInvalidSession error) *Client[T] {
	return &Client[T]{
		service:           service,
		convert:           convert,
		errInvalidSession: errInvalidSession,
	}
}

func (c *Client[T]) Whoami(ctx context.Context, sessionUUID string) (T, error) {
	var zero T

	response, err := c.service.Whoami(ctx, &iamGrpc.WhoamiRequest{
		SessionUuid: sessionUUID,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && (st.Code() == codes.NotFound || st.Code() == codes.InvalidArgument) {
			return zero, c.errInvalidSession
		}

		return zero, fmt.Errorf("iam whoami: %w", err)
	}

	return c.convert(response.GetUser()), nil
}

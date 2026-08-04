package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/Reensef/go-microservices-course/order/internal/model"
)

const bearerPrefix = "Bearer "

// UnaryClientInterceptor прокидывает токен сессии, положенный HTTP auth middleware
// в контекст запроса, в исходящую gRPC-metadata — чтобы вызываемый сервис (inventory)
// мог независимо провалидировать сессию через iam.
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if token, ok := model.TokenFromContext(ctx); ok && token != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, "authorization", bearerPrefix+token)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

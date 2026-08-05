package interceptor

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	grpcClients "github.com/Reensef/go-microservices-course/payment/internal/client/grpc"
	"github.com/Reensef/go-microservices-course/payment/internal/model"
)

const bearerPrefix = "Bearer "

var exemptMethods = map[string]struct{}{
	"/grpc.health.v1.Health/Check":                                   {},
	"/grpc.health.v1.Health/Watch":                                   {},
	"/grpc.reflection.v1.ServerReflection/ServerReflectionInfo":      {},
	"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo": {},
}

func NewAuthInterceptor(iamClient grpcClients.IAMClient) grpc.UnaryServerInterceptor {
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
			return nil, status.Error(codes.Unauthenticated, "missing or malformed authorization metadata")
		}

		token := strings.TrimPrefix(values[0], bearerPrefix)
		if token == "" {
			return nil, status.Error(codes.Unauthenticated, "missing or malformed authorization metadata")
		}

		user, err := iamClient.Whoami(ctx, token)
		if err != nil {
			if errors.Is(err, model.ErrInvalidSession) {
				return nil, status.Error(codes.Unauthenticated, "invalid or expired session")
			}

			return nil, status.Error(codes.Internal, "failed to validate session")
		}

		return handler(model.WithUserContext(ctx, user), req)
	}
}

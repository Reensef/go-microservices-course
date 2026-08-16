package interceptor

import (
	"google.golang.org/grpc"

	grpcClients "github.com/Reensef/go-microservices-course/payment/internal/client/grpc"
	"github.com/Reensef/go-microservices-course/payment/internal/model"
	"github.com/Reensef/go-microservices-course/platform/pkg/grpc/authinterceptor"
)

func NewAuthInterceptor(iamClient grpcClients.IAMClient) grpc.UnaryServerInterceptor {
	return authinterceptor.New(iamClient, model.WithUserContext, model.ErrInvalidSession)
}

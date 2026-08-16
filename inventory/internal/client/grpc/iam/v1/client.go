package v1

import (
	grpcClients "github.com/Reensef/go-microservices-course/inventory/internal/client/grpc"
	converter "github.com/Reensef/go-microservices-course/inventory/internal/client/grpc/iam/v1/converter"
	"github.com/Reensef/go-microservices-course/inventory/internal/model"
	"github.com/Reensef/go-microservices-course/platform/pkg/grpc/iamclient"
	iamGrpc "github.com/Reensef/go-microservices-course/shared/pkg/proto/iam/v1"
)

var _ grpcClients.IAMClient = (*iamclient.Client[model.User])(nil)

func New(service iamGrpc.AuthServiceClient) *iamclient.Client[model.User] {
	return iamclient.New(service, converter.ToModelUser, model.ErrInvalidSession)
}

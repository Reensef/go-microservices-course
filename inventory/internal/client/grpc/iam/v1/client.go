package v1

import (
	grpcClients "github.com/Reensef/go-microservices-course/inventory/internal/client/grpc"
	iamGrpc "github.com/Reensef/go-microservices-course/shared/pkg/proto/iam/v1"
)

var _ grpcClients.IAMClient = (*iamClient)(nil)

type iamClient struct {
	service iamGrpc.AuthServiceClient
}

func New(service iamGrpc.AuthServiceClient) *iamClient {
	return &iamClient{
		service: service,
	}
}

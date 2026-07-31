package v1

import (
	service "github.com/Reensef/go-microservices-course/iam/internal/service"
	iamV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/iam/v1"
)

type api struct {
	iamV1.UnimplementedUserServiceServer

	service service.UserService
}

func New(service service.UserService) *api {
	return &api{
		service: service,
	}
}

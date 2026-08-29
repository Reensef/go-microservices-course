package v1

import (
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"

	service "github.com/Reensef/go-microservices-course/iam/internal/service"
)

type api struct {
	authv3.UnimplementedAuthorizationServer

	service service.AuthService
}

func New(service service.AuthService) *api {
	return &api{
		service: service,
	}
}

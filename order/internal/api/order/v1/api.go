package v1

import (
	"context"
	"net/http"

	"github.com/Reensef/go-microservices-course/order/internal/service"
	orderV1 "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
)

var _ orderV1.Handler = (*api)(nil)

type api struct {
	orderService service.OrderService
}

func NewAPI(service service.OrderService) *api {
	return &api{
		orderService: service,
	}
}

func (a *api) NewError(
	_ context.Context,
	err error,
) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		},
	}
}

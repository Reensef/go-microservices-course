package v1

import (
	"context"
	"net/http"

	"github.com/Reensef/go-microservices-course/order/internal/service"
	orderApi "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
)

var _ orderApi.Handler = (*handler)(nil)

type handler struct {
	orderService service.OrderService
}

func NewHandler(service service.OrderService) *handler {
	return &handler{
		orderService: service,
	}
}

func (a *handler) NewError(
	_ context.Context,
	err error,
) *orderApi.GenericErrorStatusCode {
	return &orderApi.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderApi.GenericError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		},
	}
}

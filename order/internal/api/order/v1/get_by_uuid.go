package v1

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Reensef/go-microservices-course/order/internal/api/order/v1/converter"
	"github.com/Reensef/go-microservices-course/order/internal/model"
	orderApi "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
)

func (a *handler) GetOrderByUUID(
	ctx context.Context,
	params orderApi.GetOrderByUUIDParams,
) (orderApi.GetOrderByUUIDRes, error) {
	order, err := a.orderService.GetOrderByUUID(ctx, params.OrderUUID)
	if err != nil {
		log.Printf("api: error getting order by UUID: %s", err)

		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return &orderApi.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("order with UUID '%s' not found", params.OrderUUID),
			}, nil
		case errors.Is(err, model.ErrOrderUuidInvalidFormat):
			return &orderApi.ValidationError{
				Code:    422,
				Message: "order must be UUID format",
			}, nil
		default:
			return &orderApi.InternalServerError{
				Code:    500,
				Message: "internal server error",
			}, nil
		}
	}

	return converter.ToApiOrder(order), nil
}

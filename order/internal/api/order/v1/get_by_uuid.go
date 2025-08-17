package v1

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Reensef/go-microservices-course/order/internal/api/order/v1/converter"
	"github.com/Reensef/go-microservices-course/order/internal/model"
	orderV1 "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrderByUUID(
	ctx context.Context,
	params orderV1.GetOrderByUUIDParams,
) (orderV1.GetOrderByUUIDRes, error) {
	order, err := a.orderService.GetOrderByUUID(ctx, params.OrderUUID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("order with UUID '%s' not found", params.OrderUUID),
			}, nil
		case errors.Is(err, model.ErrOrderUuidInvalidFormat):
			return &orderV1.ValidationError{
				Code:    422,
				Message: "order must be UUID format",
			}, nil
		default:
			log.Printf("api: error getting order by UUID: %s", err)

			return &orderV1.InternalServerError{
				Code:    500,
				Message: "internal server error",
			}, nil
		}
	}

	return converter.ToApiOrder(order), nil
}

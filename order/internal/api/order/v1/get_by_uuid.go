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
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Order with UUID '%s' not found", params.OrderUUID.String()),
			}, nil
		}

		log.Printf("api: error getting order by UUID: %s", err)
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	return converter.ToApiOrder(order), nil
}

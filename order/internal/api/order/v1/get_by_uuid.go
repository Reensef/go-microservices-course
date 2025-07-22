package v1

import (
	"context"
	"errors"

	"github.com/Reensef/go-microservices-course/order/internal/model"
	orderV1 "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrderByUUID(
	ctx context.Context,
	params orderV1.GetOrderByUUIDParams,
) (orderV1.GetOrderByUUIDRes, error) {
	order, err := a.orderService.GetOrderByUUID(ctx, params.OrderUUID)

	if errors.Is(err, model.ErrOrderNotFound) {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order by UUID '" + params.OrderUUID.String() + "' not found",
		}, nil
	} else if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal server error",
		}, nil
	} else {
		return &orderV1.Get{
			OrderUUID:  order.Uuid,
			TotalPrice: order.Info.TotalPrice,
		}, nil
	}

	conver

	order := o.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order by UUID '" + params.OrderUUID.String() + "' not found",
		}, nil
	}

	return order, nil
}

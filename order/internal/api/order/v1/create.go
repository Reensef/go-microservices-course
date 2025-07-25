package v1

import (
	"context"
	"errors"

	"github.com/Reensef/go-microservices-course/order/internal/model"
	orderV1 "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(
	ctx context.Context,
	req *orderV1.CreateOrderRequest,
) (orderV1.CreateOrderRes, error) {
	order, err := a.orderService.CreateOrder(ctx, req.GetUserUUID(), req.GetPartUuids())

	if errors.Is(err, model.ErrPartNotFound) {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "part not found",
		}, nil
	} else if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "internal server error",
		}, nil
	} else {
		return &orderV1.CreateOrderResponse{
			OrderUUID:  order.Uuid,
			TotalPrice: order.Info.TotalPrice,
		}, nil
	}
}

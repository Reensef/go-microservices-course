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
	orderInfo := &model.OrderInfo{
		UserUuid:  req.UserUUID,
		PartUuids: req.PartUuids,
	}
	orderUuid, err := a.orderService.CreateOrder(
		ctx,
		orderInfo,
	)

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
			OrderUUID: *orderUuid,
			// Question: Стоит ли обновлять OrderInfo внутри CreateOrder
			// или стоит возвращать totalPrice явно?
			TotalPrice: orderInfo.TotalPrice,
		}, nil
	}
}

package v1

import (
	"context"
	"errors"
	"log"

	"github.com/Reensef/go-microservices-course/order/internal/model"
	orderV1 "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(
	ctx context.Context,
	req *orderV1.CreateOrderRequest,
) (orderV1.CreateOrderRes, error) {
	orderInfo := &model.OrderInfo{
		UserUuid: req.GetUserUUID(),
		PartIds:  req.GetPartIds(),
	}
	order, err := a.orderService.CreateOrder(
		ctx,
		orderInfo,
	)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return &orderV1.NotFoundError{
				Code:    404,
				Message: "order not found",
			}, nil
		case errors.Is(err, model.ErrUserUuidInvalidFormat):
			return &orderV1.ValidationError{
				Code:    422,
				Message: "user UUID must be UUID format",
			}, nil
		case errors.Is(err, model.ErrPartIdInvalidFormat):
			return &orderV1.ValidationError{
				Code:    422,
				Message: "part ID must be ObjectID format",
			}, nil
		default:
			log.Printf("api: error creating order: %s", err)

			return &orderV1.InternalServerError{
				Code:    500,
				Message: "internal server error",
			}, nil
		}
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  order.Uuid,
		TotalPrice: orderInfo.TotalPrice,
	}, nil
}

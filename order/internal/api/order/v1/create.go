package v1

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"

	"github.com/Reensef/go-microservices-course/order/internal/model"
	orderV1 "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(
	ctx context.Context,
	req *orderV1.CreateOrderRequest,
) (orderV1.CreateOrderRes, error) {
	validateResponse := validateCreateOrderRequest(*req)
	if validateResponse != nil {
		return validateResponse, nil
	}

	orderInfo := &model.OrderInfo{
		UserUuid: req.GetUserUUID(),
		PartIds:  req.GetPartIds(),
	}
	order, err := a.orderService.CreateOrder(
		ctx,
		orderInfo,
	)
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: "part not found",
			}, nil
		}
		log.Printf("api: error creating order: %s", err)
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "internal server error",
		}, nil
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  order.Uuid,
		TotalPrice: orderInfo.TotalPrice,
	}, nil
}

func validateCreateOrderRequest(req orderV1.CreateOrderRequest) orderV1.CreateOrderRes {
	if uuid.Validate(req.GetUserUUID()) != nil {
		return &orderV1.ValidationError{
			Code:    422,
			Message: "User UUID must be UUID format",
		}
	}
	return nil
}

package v1

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/Reensef/go-microservices-course/order/internal/model"
	orderApi "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
	"github.com/Reensef/go-microservices-course/platform/pkg/logger"
)

func (a *handler) CreateOrder(
	ctx context.Context,
	req *orderApi.CreateOrderRequest,
) (orderApi.CreateOrderRes, error) {
	user, ok := model.UserFromContext(ctx)
	if !ok {
		return &orderApi.ForbiddenError{
			Code:    403,
			Message: "access denied",
		}, nil
	}

	orderInfo := &model.OrderInfo{
		UserUuid: user.Uuid,
		PartIds:  req.GetPartIds(),
	}
	order, err := a.orderService.CreateOrder(
		ctx,
		orderInfo,
	)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return &orderApi.NotFoundError{
				Code:    404,
				Message: "order not found",
			}, nil
		case errors.Is(err, model.ErrUserUuidInvalidFormat):
			return &orderApi.ValidationError{
				Code:    422,
				Message: "user UUID must be UUID format",
			}, nil
		case errors.Is(err, model.ErrPartIdInvalidFormat):
			return &orderApi.ValidationError{
				Code:    422,
				Message: "part ID must be ObjectID format",
			}, nil
		default:
			logger.Error("api: error creating order", zap.Error(err))

			return &orderApi.InternalServerError{
				Code:    500,
				Message: "internal server error",
			}, nil
		}
	}

	return &orderApi.CreateOrderResponse{
		OrderUUID:  order.Uuid,
		TotalPrice: orderInfo.TotalPrice,
	}, nil
}

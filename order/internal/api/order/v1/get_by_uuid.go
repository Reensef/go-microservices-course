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
	user, ok := model.UserFromContext(ctx)
	if !ok {
		return &orderApi.ForbiddenError{
			Code:    403,
			Message: "access denied",
		}, nil
	}

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

	if order.Info.UserUuid != user.Uuid {
		return &orderApi.ForbiddenError{
			Code:    403,
			Message: fmt.Sprintf("order with UUID '%s' does not belong to the authenticated user", params.OrderUUID),
		}, nil
	}

	return converter.ToApiOrder(order), nil
}

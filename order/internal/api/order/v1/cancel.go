package v1

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Reensef/go-microservices-course/order/internal/model"
	orderV1 "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(
	ctx context.Context,
	params orderV1.CancelOrderParams,
) (orderV1.CancelOrderRes, error) {
	err := a.orderService.CancelOrder(ctx, params.OrderUUID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderUuidInvalidFormat):
			return &orderV1.ValidationError{
				Code:    422,
				Message: "order must be UUID format",
			}, nil
		case errors.Is(err, model.ErrOrderNotFound):
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("order by UUID '%s' not found", params.OrderUUID),
			}, nil
		case errors.Is(err, model.ErrOrderAlreadyPaid):
			return &orderV1.ConflictError{
				Code:    409,
				Message: fmt.Sprintf("order with UUID '%s' already paid", params.OrderUUID),
			}, nil
		default:
			log.Printf("api: error cancelling order: %s", err)

			return &orderV1.InternalServerError{
				Code:    500,
				Message: "internal server error",
			}, nil
		}
	}

	return &orderV1.CancelOrderNoContent{}, nil
}

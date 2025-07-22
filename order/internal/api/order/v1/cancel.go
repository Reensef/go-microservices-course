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
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Order by UUID '%s' not found", params.OrderUUID.String()),
			}, nil
		} else if errors.Is(err, model.ErrOrderAlreadyPaid) {
			return &orderV1.ConflictError{
				Code:    409,
				Message: fmt.Sprintf("Order with UUID '%s' already paid", params.OrderUUID.String()),
			}, nil
		}

		log.Printf("api: error cancelling order: %s", err)
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	return &orderV1.CancelOrderNoContent{}, nil
}

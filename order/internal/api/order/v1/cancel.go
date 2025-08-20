package v1

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/Reensef/go-microservices-course/order/internal/model"
	orderV1 "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(
	ctx context.Context,
	params orderV1.CancelOrderParams,
) (orderV1.CancelOrderRes, error) {
	validateResult := validateCancelOrderParams(&params)
	if validateResult != nil {
		return validateResult, nil
	}

	err := a.orderService.CancelOrder(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Order by UUID '%s' not found", params.OrderUUID),
			}, nil
		} else if errors.Is(err, model.ErrOrderAlreadyPaid) {
			return &orderV1.ConflictError{
				Code:    409,
				Message: fmt.Sprintf("Order with UUID '%s' already paid", params.OrderUUID),
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

func validateCancelOrderParams(params *orderV1.CancelOrderParams) orderV1.CancelOrderRes {
	if uuid.Validate(params.OrderUUID) != nil {
		return &orderV1.ValidationError{
			Code:    422,
			Message: "Order must be UUID format",
		}
	}
	return nil
}

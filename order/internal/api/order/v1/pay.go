package v1

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/Reensef/go-microservices-course/order/internal/api/order/v1/converter"
	"github.com/Reensef/go-microservices-course/order/internal/model"
	orderV1 "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(
	ctx context.Context,
	req *orderV1.PayOrderRequest,
	params orderV1.PayOrderParams,
) (orderV1.PayOrderRes, error) {
	userUuid := uuid.NewString()

	transactionUUID, err := a.orderService.PayOrder(
		ctx,
		params.OrderUUID,
		userUuid,
		converter.ToModelPaymentMethod(req.GetPaymentMethod()),
	)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("order with UUID '%s' not found", params.OrderUUID),
			}, nil
		case errors.Is(err, model.ErrOrderAlreadyPaid):
			return &orderV1.ConflictError{
				Code:    409,
				Message: fmt.Sprintf("order with UUID '%s' already paid", params.OrderUUID),
			}, nil
		case errors.Is(err, model.ErrPaymentMethodUnspecified):
			return &orderV1.ValidationError{
				Code:    422,
				Message: "payment method must be specified",
			}, nil
		case errors.Is(err, model.ErrOrderUuidInvalidFormat):
			return &orderV1.ValidationError{
				Code:    422,
				Message: "order must be UUID format",
			}, nil
		case errors.Is(err, model.ErrUserUuidInvalidFormat):
			return &orderV1.ValidationError{
				Code:    422,
				Message: "user UUID must be UUID format",
			}, nil
		default:
			log.Printf("api: error paying order: %s", err)

			return &orderV1.InternalServerError{
				Code:    500,
				Message: "internal server error",
			}, nil
		}
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: orderV1.NewOptString(*transactionUUID),
	}, nil
}

package v1

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/Reensef/go-microservices-course/order/internal/api/order/v1/converter"
	"github.com/Reensef/go-microservices-course/order/internal/model"
	orderApi "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
)

func (a *handler) PayOrder(
	ctx context.Context,
	req *orderApi.PayOrderRequest,
	params orderApi.PayOrderParams,
) (orderApi.PayOrderRes, error) {
	userUuid := uuid.NewString()
	transactionUUID, err := a.orderService.PayOrder(
		ctx,
		params.OrderUUID,
		userUuid,
		converter.ToModelPaymentMethod(req.GetPaymentMethod()),
	)
	if err != nil {
		log.Printf("api: error paying order: %s", err)

		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return &orderApi.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("order with UUID '%s' not found", params.OrderUUID),
			}, nil
		case errors.Is(err, model.ErrOrderAlreadyPaid):
			return &orderApi.ConflictError{
				Code:    409,
				Message: fmt.Sprintf("order with UUID '%s' already paid", params.OrderUUID),
			}, nil
		case errors.Is(err, model.ErrPaymentMethodUnspecified):
			return &orderApi.ValidationError{
				Code:    422,
				Message: "payment method must be specified",
			}, nil
		case errors.Is(err, model.ErrOrderUuidInvalidFormat):
			return &orderApi.ValidationError{
				Code:    422,
				Message: "order must be UUID format",
			}, nil
		case errors.Is(err, model.ErrUserUuidInvalidFormat):
			return &orderApi.ValidationError{
				Code:    422,
				Message: "user UUID must be UUID format",
			}, nil
		default:
			return &orderApi.InternalServerError{
				Code:    500,
				Message: "internal server error",
			}, nil
		}
	}

	return &orderApi.PayOrderResponse{
		TransactionUUID: orderApi.NewOptString(*transactionUUID),
	}, nil
}

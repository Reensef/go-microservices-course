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
	userUuid := uuid.New()

	transactionUUID, err := a.orderService.PayOrder(
		ctx,
		params.OrderUUID,
		userUuid,
		converter.ToModelPaymentMethod(req.GetPaymentMethod()),
	)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Order with UUID '%s' not found", params.OrderUUID.String()),
			}, nil
		} else if errors.Is(err, model.ErrOrderAlreadyPaid) {
			return &orderV1.ConflictError{
				Code:    409,
				Message: fmt.Sprintf("Order with UUID '%s' already paid", params.OrderUUID.String()),
			}, nil
		}

		log.Printf("api: error paying order: %s", err)
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: orderV1.NewOptUUID(*transactionUUID),
	}, nil
}

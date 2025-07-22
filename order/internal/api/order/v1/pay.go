package v1

import (
	"context"

	"github.com/google/uuid"
)

func (a *api) PayOrder(
	ctx context.Context,
	req *orderV1.PayOrderRequest,
	params orderV1.PayOrderParams,
) (orderV1.PayOrderRes, error) {
	order := o.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order by uuid '" + params.OrderUUID.String() + "' not found",
		}, nil
	}

	resPayOrder, err := o.paymentService.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid: order.OrderUUID.String(),
		UserUuid:  order.UserUUID.String(),
	})
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	transactionUUID, err := uuid.Parse(resPayOrder.TransactionUuid)
	if err != nil {
		return nil, err
	}

	order.TransactionUUID.SetTo(transactionUUID)
	order.PaymentMethod.SetTo(req.PaymentMethod)
	order.Status = orderV1.OrderStatusPAID

	err = o.storage.UpdateOrder(order)
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order by uuid '" + params.OrderUUID.String() + "' not found",
		}, nil
	}

	response := &orderV1.PayOrderResponse{}
	response.TransactionUUID.SetTo(transactionUUID)

	return response, nil
}

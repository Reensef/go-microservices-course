package v1

import (
	"context"

	"github.com/google/uuid"

	protoConverter "github.com/Reensef/go-microservices-course/order/internal/client/grpc/payment/v1/converter"
	"github.com/Reensef/go-microservices-course/order/internal/model"
	paymentV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/payment/v1"
)

func (c *paymentV1Client) PayOrder(
	ctx context.Context,
	orderUuid, userUuid uuid.UUID,
	paymentMethod model.OrderPaymentMethod,
) (*uuid.UUID, error) {
	response, err := c.service.PayOrder(
		ctx,
		&paymentV1.PayOrderRequest{
			OrderUuid:     orderUuid.String(),
			UserUuid:      userUuid.String(),
			PaymentMethod: protoConverter.ToProtoPaymentMethod(paymentMethod),
		},
	)
	if err != nil {
		return nil, err
	}

	transactionUuid, err := uuid.Parse(response.GetTransactionUuid())
	if err != nil {
		return nil, err
	}

	return &transactionUuid, nil
}

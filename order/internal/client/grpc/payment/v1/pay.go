package v1

import (
	"context"

	protoConverter "github.com/Reensef/go-microservices-course/order/internal/client/grpc/payment/v1/converter"
	"github.com/Reensef/go-microservices-course/order/internal/model"
	paymentV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/payment/v1"
)

func (c *paymentV1Client) PayOrder(
	ctx context.Context,
	orderUuid, userUuid string,
	paymentMethod model.OrderPaymentMethod,
) (*string, error) {
	response, err := c.service.PayOrder(
		ctx,
		&paymentV1.PayOrderRequest{
			OrderUuid:     orderUuid,
			UserUuid:      userUuid,
			PaymentMethod: protoConverter.ToProtoPaymentMethod(paymentMethod),
		},
	)
	if err != nil {
		return nil, err
	}

	transactionUuid := response.GetTransactionUuid()

	return &transactionUuid, nil
}

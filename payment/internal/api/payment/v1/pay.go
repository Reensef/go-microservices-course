package v1

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	converter "github.com/Reensef/go-microservices-course/payment/internal/api/payment/v1/converter"
	paymentV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(
	ctx context.Context,
	req *paymentV1.PayOrderRequest,
) (*paymentV1.PayOrderResponse, error) {
	orderUuid, err := uuid.Parse(req.OrderUuid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "")
	}

	userUuid, err := uuid.Parse(req.UserUuid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "")
	}

	transactionUuid, err := a.service.Pay(
		ctx,
		&orderUuid,
		&userUuid,
		converter.ProtoPaymentMethodToModel(req.PaymentMethod),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "")
	}

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUuid.String(),
	}, nil
}

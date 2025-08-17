package v1

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	converter "github.com/Reensef/go-microservices-course/payment/internal/api/payment/v1/converter"
	"github.com/Reensef/go-microservices-course/payment/internal/model"
	paymentV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(
	ctx context.Context,
	req *paymentV1.PayOrderRequest,
) (*paymentV1.PayOrderResponse, error) {
	orderUuid := req.GetOrderUuid()
	userUuid := req.GetUserUuid()

	transactionUuid, err := a.service.Pay(
		ctx,
		orderUuid,
		userUuid,
		converter.ToModelPaymentMethod(req.GetPaymentMethod()),
	)
	if err != nil {
		log.Printf("api: error paying order: %s", err.Error())

		switch {
		case errors.Is(err, model.ErrOrderUuidInvalidFormat):
			return nil, status.Errorf(codes.InvalidArgument, "order must be UUID format")
		case errors.Is(err, model.ErrUserUuidInvalidFormat):
			return nil, status.Errorf(codes.InvalidArgument, "user UUID must be UUID format")
		case errors.Is(err, model.ErrPaymentMethodUnspecified):
			return nil, status.Errorf(codes.InvalidArgument, "payment method must be specified")
		default:
			return nil, status.Errorf(codes.Internal, "internal server error")
		}
	}

	log.Println("Оплата прошла успешно, transaction_uuid:", transactionUuid)

	return &paymentV1.PayOrderResponse{
		TransactionUuid: *transactionUuid,
	}, nil
}

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
		if errors.Is(err, model.ErrPaymentMethodUnspecified) {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
		}

		log.Printf("api: error paying order: %s", err.Error())
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &paymentV1.PayOrderResponse{
		TransactionUuid: *transactionUuid,
	}, nil
}

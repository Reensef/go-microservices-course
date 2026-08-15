package v1

import (
	"context"
	"errors"
	"log"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	converter "github.com/Reensef/go-microservices-course/payment/internal/api/payment/v1/converter"
	"github.com/Reensef/go-microservices-course/payment/internal/model"
	"github.com/Reensef/go-microservices-course/platform/pkg/logger"
	paymentV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(
	ctx context.Context,
	req *paymentV1.PayOrderRequest,
) (*paymentV1.PayOrderResponse, error) {
	userUuid := req.GetUserUuid()
	authenticatedUser, ok := model.UserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	if authenticatedUser.Uuid != userUuid {
		return nil, status.Errorf(
			codes.PermissionDenied, "authenticated user does not match user_uuid in request",
		)
	}

	paymentMethod := converter.ToModelPaymentMethod(req.GetPaymentMethod())
	orderUuid := req.GetOrderUuid()

	transactionUuid, err := a.service.Pay(ctx, orderUuid, userUuid, paymentMethod)
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

	logger.Info("payment succeeded",
		zap.String("order_uuid", orderUuid),
		zap.String("transaction_uuid", *transactionUuid),
	)

	return &paymentV1.PayOrderResponse{
		TransactionUuid: *transactionUuid,
	}, nil
}

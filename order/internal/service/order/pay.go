package order

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Reensef/go-microservices-course/order/internal/model"
	"github.com/Reensef/go-microservices-course/platform/pkg/logger"
	"github.com/Reensef/go-microservices-course/platform/pkg/tracer"
)

func (s *service) PayOrder(
	ctx context.Context,
	orderUuid string,
	userUuid string,
	paymentMethod model.OrderPaymentMethod,
) (string, error) {
	ctx, span := tracer.StartSpan(ctx, "order.pay_order")
	defer span.End()

	if uuid.Validate(orderUuid) != nil {
		tracer.RecordError(span, model.ErrOrderUuidInvalidFormat)
		return "", model.ErrOrderUuidInvalidFormat
	}
	if uuid.Validate(userUuid) != nil {
		tracer.RecordError(span, model.ErrUserUuidInvalidFormat)
		return "", model.ErrUserUuidInvalidFormat
	}
	if paymentMethod == model.OrderPaymentMethod_UNSPECIFIED {
		tracer.RecordError(span, model.ErrPaymentMethodUnspecified)
		return "", model.ErrPaymentMethodUnspecified
	}

	logger.Info("paying order",
		zap.String("order_uuid", orderUuid),
		zap.String("user_uuid", userUuid),
	)

	order, err := s.orderRepo.GetOrderByUUID(ctx, orderUuid)
	if err != nil {
		tracer.RecordError(span, err)
		return "", err
	}

	if order.Info.UserUuid != userUuid {
		tracer.RecordError(span, model.ErrOrderAccessDenied)
		return "", model.ErrOrderAccessDenied
	}

	if order.Info.Status == model.OrderStatus_PAID {
		tracer.RecordError(span, model.ErrOrderAlreadyPaid)
		return "", model.ErrOrderAlreadyPaid
	}

	transactionUuid, err := s.paymentService.PayOrder(ctx, orderUuid, userUuid, paymentMethod)
	if err != nil {
		logger.Error("payment failed", zap.String("order_uuid", orderUuid), zap.Error(err))
		tracer.RecordError(span, err)
		return "", err
	}

	err = s.orderRepo.PayOrder(ctx, orderUuid, transactionUuid, paymentMethod)
	if err != nil {
		logger.Error("failed to update order payment status", zap.String("order_uuid", orderUuid), zap.Error(err))
		tracer.RecordError(span, err)
		return "", err
	}

	err = s.orderProducer.ProduceOrderPaid(ctx, model.OrderPaidEvent{
		UUID:            uuid.New().String(),
		OrderUUID:       orderUuid,
		UserUUID:        userUuid,
		TransactionUUID: transactionUuid,
		PaymentMethod:   paymentMethod,
	})
	if err != nil {
		logger.Error("failed to produce order paid event", zap.String("order_uuid", orderUuid), zap.Error(err))
		tracer.RecordError(span, err)
		return "", err
	}

	logger.Info("order paid",
		zap.String("order_uuid", orderUuid),
		zap.String("transaction_uuid", transactionUuid),
	)

	return transactionUuid, nil
}

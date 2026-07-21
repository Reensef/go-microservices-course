package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/Reensef/go-microservices-course/order/internal/model"
)

func (s *service) PayOrder(
	ctx context.Context,
	orderUuid string,
	userUuid string,
	paymentMethod model.OrderPaymentMethod,
) (string, error) {
	if uuid.Validate(orderUuid) != nil {
		return "", model.ErrOrderUuidInvalidFormat
	}
	if uuid.Validate(userUuid) != nil {
		return "", model.ErrUserUuidInvalidFormat
	}
	if paymentMethod == model.OrderPaymentMethod_UNSPECIFIED {
		return "", model.ErrPaymentMethodUnspecified
	}

	order, err := s.orderRepo.GetOrderByUUID(ctx, orderUuid)
	if err != nil {
		return "", err
	}

	if order.Info.Status == model.OrderStatus_PAID {
		return "", model.ErrOrderAlreadyPaid
	}

	transactionUuid, err := s.paymentService.PayOrder(ctx, orderUuid, userUuid, paymentMethod)
	if err != nil {
		return "", err
	}

	err = s.orderRepo.PayOrder(ctx, orderUuid, transactionUuid, paymentMethod)
	if err != nil {
		return "", err
	}

	err = s.orderProducer.ProduceOrderPaid(ctx, model.OrderPaidEvent{
		UUID:          transactionUuid,
		OrderUUID:     orderUuid,
		UserUUID:      userUuid,
		PaymentMethod: paymentMethod,
	})
	if err != nil {
		return "", err
	}

	return transactionUuid, nil
}

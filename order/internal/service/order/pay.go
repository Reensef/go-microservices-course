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
) (*string, error) {
	if uuid.Validate(orderUuid) != nil {
		return nil, model.ErrOrderUuidInvalidFormat
	}
	if uuid.Validate(userUuid) != nil {
		return nil, model.ErrUserUuidInvalidFormat
	}
	if paymentMethod == model.OrderPaymentMethod_UNSPECIFIED {
		return nil, model.ErrPaymentMethodUnspecified
	}

	status, err := s.orderRepo.GetOrderStatus(ctx, orderUuid)
	if err != nil {
		return nil, err
	}

	if status == model.OrderStatus_PAID {
		return nil, model.ErrOrderAlreadyPaid
	}

	transactionUuid, err := s.paymentService.PayOrder(ctx, orderUuid, userUuid, paymentMethod)
	if err != nil {
		return nil, err
	}

	err = s.orderRepo.PayOrder(ctx, orderUuid, *transactionUuid, paymentMethod)
	if err != nil {
		return nil, err
	}

	return transactionUuid, nil
}

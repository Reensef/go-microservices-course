package payment

import (
	"context"

	"github.com/google/uuid"

	"github.com/Reensef/go-microservices-course/payment/internal/model"
	"github.com/Reensef/go-microservices-course/platform/pkg/tracer"
)

func (s *service) Pay(
	ctx context.Context,
	orderUuid, userUuid string,
	paymentMethod model.PaymentMethod,
) (*string, error) {
	ctx, span := tracer.StartSpan(ctx, "payment.pay")
	defer span.End()

	if uuid.Validate(orderUuid) != nil {
		tracer.RecordError(span, model.ErrOrderUuidInvalidFormat)
		return nil, model.ErrOrderUuidInvalidFormat
	}
	if uuid.Validate(userUuid) != nil {
		tracer.RecordError(span, model.ErrUserUuidInvalidFormat)
		return nil, model.ErrUserUuidInvalidFormat
	}
	if paymentMethod == model.PaymentMethod_UNSPECIFIED {
		tracer.RecordError(span, model.ErrPaymentMethodUnspecified)
		return nil, model.ErrPaymentMethodUnspecified
	}

	transactionUuid := uuid.NewString()

	return &transactionUuid, nil
}

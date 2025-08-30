package payment

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Reensef/go-microservices-course/payment/internal/model"
)

func TestPay(t *testing.T) {
	t.Run("Unspecified payment method", func(t *testing.T) {
		s := &service{}
		orderUuid := uuid.NewString()
		userUuid := uuid.NewString()
		paymentMethod := model.PaymentMethod_UNSPECIFIED
		_, err := s.Pay(context.Background(), orderUuid, userUuid, paymentMethod)
		assert.EqualError(t, err, model.ErrPaymentMethodUnspecified.Error())
	})
	t.Run("Successful payment", func(t *testing.T) {
		s := &service{}
		orderUuid := uuid.NewString()
		userUuid := uuid.NewString()
		paymentMethod := model.PaymentMethod_CARD
		transactionUuid, err := s.Pay(context.Background(), orderUuid, userUuid, paymentMethod)
		assert.NoError(t, err)
		assert.NotNil(t, transactionUuid)
	})
}

package v1

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	converter "github.com/Reensef/go-microservices-course/payment/internal/api/payment/v1/converter"
	"github.com/Reensef/go-microservices-course/payment/internal/model"
	"github.com/Reensef/go-microservices-course/payment/internal/service/mocks"
	paymentV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/payment/v1"
)

func TestPayOrder(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockPaymentService := mocks.NewMockPaymentService(t)
		apiInstance := New(mockPaymentService)

		orderUuid := "550e8400-e29b-41d4-a716-446655440000"
		userUuid := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
		paymentMethod := paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD

		ctx := model.WithUserContext(context.Background(), model.User{Uuid: userUuid})

		req := &paymentV1.PayOrderRequest{
			OrderUuid:     orderUuid,
			UserUuid:      userUuid,
			PaymentMethod: paymentMethod,
		}

		expectedTransactionUuid := uuid.NewString()

		mockPaymentService.EXPECT().
			Pay(ctx, orderUuid, userUuid, converter.ToModelPaymentMethod(paymentMethod)).
			Return(&expectedTransactionUuid, nil)

		resp, err := apiInstance.PayOrder(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, expectedTransactionUuid, resp.TransactionUuid)
	})

	t.Run("Unspecified payment method", func(t *testing.T) {
		mockPaymentService := mocks.NewMockPaymentService(t)
		apiInstance := New(mockPaymentService)

		orderUuid := "550e8400-e29b-41d4-a716-446655440000"
		userUuid := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"

		ctx := model.WithUserContext(context.Background(), model.User{Uuid: userUuid})

		req := &paymentV1.PayOrderRequest{
			OrderUuid:     orderUuid,
			UserUuid:      userUuid,
			PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED,
		}

		mockPaymentService.EXPECT().
			Pay(ctx, orderUuid, userUuid, converter.ToModelPaymentMethod(req.PaymentMethod)).
			Return(nil, model.ErrPaymentMethodUnspecified)

		resp, err := apiInstance.PayOrder(ctx, req)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("Internal error", func(t *testing.T) {
		mockPaymentService := mocks.NewMockPaymentService(t)
		apiInstance := New(mockPaymentService)

		orderUuid := "550e8400-e29b-41d4-a716-446655440000"
		userUuid := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
		paymentMethod := paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD

		ctx := model.WithUserContext(context.Background(), model.User{Uuid: userUuid})

		req := &paymentV1.PayOrderRequest{
			OrderUuid:     orderUuid,
			UserUuid:      userUuid,
			PaymentMethod: paymentMethod,
		}

		mockPaymentService.EXPECT().
			Pay(ctx, orderUuid, userUuid, converter.ToModelPaymentMethod(paymentMethod)).
			Return(nil, assert.AnError)

		resp, err := apiInstance.PayOrder(ctx, req)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.Equal(t, codes.Internal, status.Code(err))
	})

	t.Run("No authenticated user in context", func(t *testing.T) {
		mockPaymentService := mocks.NewMockPaymentService(t)
		apiInstance := New(mockPaymentService)

		req := &paymentV1.PayOrderRequest{
			OrderUuid:     "550e8400-e29b-41d4-a716-446655440000",
			UserUuid:      "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		}

		resp, err := apiInstance.PayOrder(context.Background(), req)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.Equal(t, codes.Internal, status.Code(err))
		assert.Empty(t, mockPaymentService.Calls)
	})

	t.Run("Authenticated user does not match request user_uuid", func(t *testing.T) {
		mockPaymentService := mocks.NewMockPaymentService(t)
		apiInstance := New(mockPaymentService)

		userUuid := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
		otherUuid := "11111111-1111-1111-1111-111111111111"

		ctx := model.WithUserContext(context.Background(), model.User{Uuid: otherUuid})

		req := &paymentV1.PayOrderRequest{
			OrderUuid:     "550e8400-e29b-41d4-a716-446655440000",
			UserUuid:      userUuid,
			PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		}

		resp, err := apiInstance.PayOrder(ctx, req)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.Equal(t, codes.PermissionDenied, status.Code(err))
		assert.Empty(t, mockPaymentService.Calls)
	})
}

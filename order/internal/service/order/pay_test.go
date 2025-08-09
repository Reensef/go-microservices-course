package order

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	grpcMocks "github.com/Reensef/go-microservices-course/order/internal/client/grpc/mocks"
	"github.com/Reensef/go-microservices-course/order/internal/model"
	repoMocks "github.com/Reensef/go-microservices-course/order/internal/repository/mocks"
)

func TestPayOrder(t *testing.T) {
	t.Run("Error from payment service", func(t *testing.T) {
		repo := repoMocks.NewMockOrderRepository(t)
		inventory := grpcMocks.NewMockIntentoryServiceClient(t)
		payment := grpcMocks.NewMockPaymentServiceClient(t)
		service := NewService(repo, inventory, payment)

		userUuid := uuid.New()
		orderUuid := uuid.New()
		paymentMethod := model.OrderPaymentMethod_CARD
		paymentError := fmt.Errorf("error")

		payment.EXPECT().PayOrder(context.Background(), orderUuid, userUuid, paymentMethod).Return(nil, paymentError).Once()

		uuid, err := service.PayOrder(context.Background(), orderUuid, userUuid, paymentMethod)

		assert.Nil(t, uuid)
		assert.Equal(t, err, paymentError)

		assert.Empty(t, repo.Calls)
		assert.Empty(t, inventory.Calls)
	})

	t.Run("Error from repository", func(t *testing.T) {
		repo := repoMocks.NewMockOrderRepository(t)
		inventory := grpcMocks.NewMockIntentoryServiceClient(t)
		payment := grpcMocks.NewMockPaymentServiceClient(t)
		service := NewService(repo, inventory, payment)

		userUuid := uuid.New()
		orderUuid := uuid.New()
		paymentMethod := model.OrderPaymentMethod_CARD
		repoError := fmt.Errorf("error")
		transactionUuid := uuid.New()

		payment.EXPECT().PayOrder(context.Background(), orderUuid, userUuid, paymentMethod).Return(&transactionUuid, nil).Once()
		repo.EXPECT().PayOrder(context.Background(), orderUuid, transactionUuid, paymentMethod).Return(repoError).Once()

		uuid, err := service.PayOrder(context.Background(), orderUuid, userUuid, paymentMethod)

		assert.Nil(t, uuid)
		assert.Equal(t, err, repoError)

		assert.Empty(t, inventory.Calls)
	})

	t.Run("Success", func(t *testing.T) {
		repo := repoMocks.NewMockOrderRepository(t)
		inventory := grpcMocks.NewMockIntentoryServiceClient(t)
		payment := grpcMocks.NewMockPaymentServiceClient(t)
		service := NewService(repo, inventory, payment)

		userUuid := uuid.New()
		orderUuid := uuid.New()
		paymentMethod := model.OrderPaymentMethod_CARD
		transactionUuid := uuid.New()

		payment.EXPECT().PayOrder(context.Background(), orderUuid, userUuid, paymentMethod).Return(&transactionUuid, nil).Once()
		repo.EXPECT().PayOrder(context.Background(), orderUuid, transactionUuid, paymentMethod).Return(nil).Once()

		uuid, err := service.PayOrder(context.Background(), orderUuid, userUuid, paymentMethod)

		assert.NoError(t, err)
		assert.Equal(t, uuid, &transactionUuid)

		assert.Empty(t, inventory.Calls)
	})
}

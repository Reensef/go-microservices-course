package order

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	model "github.com/Reensef/go-microservices-course/order/internal/model"
	mocks "github.com/Reensef/go-microservices-course/order/internal/repository/mocks"
)

func TestCancelOrder(t *testing.T) {
	t.Run("Order not found", func(t *testing.T) {
		repo := mocks.NewMockOrderRepository(t)
		service := New(repo, nil, nil, nil)

		orderUuid := uuid.NewString()
		requesterUuid := uuid.NewString()
		repo.EXPECT().GetOrderByUUID(context.Background(), orderUuid).
			Return(nil, model.ErrOrderNotFound).
			Once()

		err := service.CancelOrder(context.Background(), orderUuid, requesterUuid)

		assert.Equal(t, model.ErrOrderNotFound, err)
	})

	t.Run("Order found", func(t *testing.T) {
		repo := mocks.NewMockOrderRepository(t)
		service := New(repo, nil, nil, nil)

		orderUuid := uuid.NewString()
		requesterUuid := uuid.NewString()
		order := &model.Order{Info: model.OrderInfo{UserUuid: requesterUuid}}
		repo.EXPECT().GetOrderByUUID(context.Background(), orderUuid).
			Return(order, nil).
			Once()

		repo.EXPECT().CancelOrder(context.Background(), orderUuid).
			Return(nil).Once()

		err := service.CancelOrder(context.Background(), orderUuid, requesterUuid)

		assert.NoError(t, err)
	})

	t.Run("Access denied", func(t *testing.T) {
		repo := mocks.NewMockOrderRepository(t)
		service := New(repo, nil, nil, nil)

		orderUuid := uuid.NewString()
		requesterUuid := uuid.NewString()
		ownerUuid := uuid.NewString()
		order := &model.Order{Info: model.OrderInfo{UserUuid: ownerUuid}}
		repo.EXPECT().GetOrderByUUID(context.Background(), orderUuid).
			Return(order, nil).
			Once()

		err := service.CancelOrder(context.Background(), orderUuid, requesterUuid)

		assert.ErrorIs(t, err, model.ErrOrderAccessDenied)
	})
}

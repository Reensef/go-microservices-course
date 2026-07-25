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

		uuid := uuid.NewString()
		repo.EXPECT().GetOrderByUUID(context.Background(), uuid).
			Return(nil, model.ErrOrderNotFound).
			Once()

		err := service.CancelOrder(context.Background(), uuid)

		assert.Equal(t, model.ErrOrderNotFound, err)
	})

	t.Run("Order found", func(t *testing.T) {
		repo := mocks.NewMockOrderRepository(t)
		service := New(repo, nil, nil, nil)

		uuid := uuid.NewString()
		order := &model.Order{}
		repo.EXPECT().GetOrderByUUID(context.Background(), uuid).
			Return(order, nil).
			Once()

		repo.EXPECT().CancelOrder(context.Background(), uuid).
			Return(nil).Once()

		err := service.CancelOrder(context.Background(), uuid)

		assert.NoError(t, err)
	})
}

package order

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Reensef/go-microservices-course/order/internal/model"
	mocks "github.com/Reensef/go-microservices-course/order/internal/repository/mocks"
)

func TestCancelOrder(t *testing.T) {
	t.Run("Order not found", func(t *testing.T) {
		repo := mocks.NewMockOrderRepository(t)
		service := NewService(repo, nil, nil)

		uuid := uuid.New()
		repo.EXPECT().CancelOrder(context.Background(), uuid).Return(model.ErrOrderNotFound).Once()

		err := service.CancelOrder(context.Background(), uuid)

		assert.Equal(t, model.ErrOrderNotFound, err)
	})

	t.Run("Order found", func(t *testing.T) {
		repo := mocks.NewMockOrderRepository(t)
		service := NewService(repo, nil, nil)

		uuid := uuid.New()
		repo.EXPECT().CancelOrder(context.Background(), uuid).Return(nil).Once()

		err := service.CancelOrder(context.Background(), uuid)

		assert.NoError(t, err)
	})
}

package order

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Reensef/go-microservices-course/order/internal/model"
	repoModel "github.com/Reensef/go-microservices-course/order/internal/repository/model"
)

func TestCancelOrder(t *testing.T) {
	t.Run("Test canceling a non-existent order", func(t *testing.T) {
		repo := &repository{
			data: make(map[uuid.UUID]*repoModel.Order),
		}
		orderUuid := uuid.New()
		err := repo.CancelOrder(context.Background(), orderUuid)
		assert.Equal(t, model.ErrOrderNotFound, err)
	})

	t.Run("Test canceling an order that has already been paid", func(t *testing.T) {
		orderUuid := uuid.New()

		repo := &repository{
			data: map[uuid.UUID]*repoModel.Order{
				orderUuid: {
					Info: repoModel.OrderInfo{
						Status: repoModel.OrderStatus_PAID,
					},
				},
			},
		}
		err := repo.CancelOrder(context.Background(), orderUuid)
		assert.Equal(t, model.ErrOrderAlreadyPaid, err)
	})

	t.Run("Test canceling a pending order successfully", func(t *testing.T) {
		orderUuid := uuid.New()

		repo := &repository{
			data: map[uuid.UUID]*repoModel.Order{
				orderUuid: {
					Info: repoModel.OrderInfo{
						Status: repoModel.OrderStatus_PENDING_PAYMENT,
					},
				},
			},
		}

		err := repo.CancelOrder(context.Background(), orderUuid)
		assert.NoError(t, err)
		order, ok := repo.data[orderUuid]
		require.True(t, ok)
		assert.Equal(t, repoModel.OrderStatus(repoModel.OrderStatus_CANCELED), order.Info.Status)
	})
}

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

func TestGetOrderByUUID(t *testing.T) {
	t.Run("Order found", func(t *testing.T) {
		orderUuid := uuid.New()
		order := &repoModel.Order{
			Uuid: orderUuid,
			Info: repoModel.OrderInfo{},
		}
		repo := &repository{
			data: map[uuid.UUID]*repoModel.Order{
				orderUuid: order,
			},
		}

		actualOrder, err := repo.GetOrderByUUID(context.Background(), orderUuid)

		require.NoError(t, err)
		assert.NotNil(t, actualOrder)
		assert.Equal(t, orderUuid, actualOrder.Uuid)
	})

	t.Run("Order not found", func(t *testing.T) {
		orderUuid := uuid.New()
		repo := &repository{
			data: map[uuid.UUID]*repoModel.Order{},
		}

		actualOrder, err := repo.GetOrderByUUID(context.Background(), orderUuid)

		require.Equal(t, model.ErrOrderNotFound, err)
		assert.Nil(t, actualOrder)
	})

	t.Run("Empty repository", func(t *testing.T) {
		repo := &repository{
			data: map[uuid.UUID]*repoModel.Order{},
		}

		uuid := uuid.New()
		actualOrder, err := repo.GetOrderByUUID(context.Background(), uuid)

		require.Equal(t, model.ErrOrderNotFound, err)
		assert.Nil(t, actualOrder)
	})
}

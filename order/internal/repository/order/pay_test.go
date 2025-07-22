package order

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Reensef/go-microservices-course/order/internal/model"
	repoModel "github.com/Reensef/go-microservices-course/order/internal/repository/model"
)

func TestPayOrder(t *testing.T) {
	t.Run("Order not found", func(t *testing.T) {
		repo := &repository{
			data: map[uuid.UUID]*repoModel.Order{},
		}

		orderUuid := uuid.New()
		transactionUUID := uuid.New()
		paymentMethod := model.OrderPaymentMethod_CARD
		expectedError := model.ErrOrderNotFound

		err := repo.PayOrder(context.Background(), orderUuid, transactionUUID, paymentMethod)

		assert.EqualError(t, err, expectedError.Error())
	})

	t.Run("Successful payment", func(t *testing.T) {
		repo := &repository{
			data: map[uuid.UUID]*repoModel.Order{},
		}

		orderUuid := uuid.New()
		transactionUUID := uuid.New()
		paymentMethod := model.OrderPaymentMethod_CARD

		repo.data[orderUuid] = &repoModel.Order{
			Info: repoModel.OrderInfo{
				Status: repoModel.OrderStatus_PENDING_PAYMENT,
			},
		}

		err := repo.PayOrder(context.Background(), orderUuid, transactionUUID, paymentMethod)

		assert.NoError(t, err)
		assert.Equal(t,
			repoModel.OrderStatus(model.OrderStatus_PAID),
			repo.data[orderUuid].Info.Status,
		)
		assert.Equal(t, transactionUUID, repo.data[orderUuid].Info.TransactionUuid)
	})
}

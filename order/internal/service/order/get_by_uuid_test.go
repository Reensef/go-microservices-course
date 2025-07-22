package order

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	grpcMocks "github.com/Reensef/go-microservices-course/order/internal/client/grpc/mocks"
	repoMocks "github.com/Reensef/go-microservices-course/order/internal/repository/mocks"
)

func TestGetOrderByUUID(t *testing.T) {
	t.Run("Error get from repo", func(t *testing.T) {
		repo := repoMocks.NewMockOrderRepository(t)
		inventory := grpcMocks.NewMockIntentoryServiceClient(t)
		payment := grpcMocks.NewMockPaymentServiceClient(t)
		service := NewService(repo, inventory, payment)

		repoError := fmt.Errorf("error")

		repo.EXPECT().GetOrderByUUID(context.Background(), mock.Anything).Return(nil, repoError).Once()

		orderUuid := uuid.New()
		order, err := service.GetOrderByUUID(context.Background(), orderUuid)

		assert.Nil(t, order)
		assert.Empty(t, inventory.Calls)
		assert.Empty(t, payment.Calls)
		assert.Equal(t, err, repoError)
	})
}

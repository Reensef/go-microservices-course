package order

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	grpcMocks "github.com/Reensef/go-microservices-course/order/internal/client/grpc/mocks"
	"github.com/Reensef/go-microservices-course/order/internal/model"
	repoMocks "github.com/Reensef/go-microservices-course/order/internal/repository/mocks"
)

func TestCreateOrder(t *testing.T) {
	t.Run("Nil info", func(t *testing.T) {
		repo := repoMocks.NewMockOrderRepository(t)
		inventory := grpcMocks.NewMockIntentoryServiceClient(t)
		payment := grpcMocks.NewMockPaymentServiceClient(t)
		service := New(repo, inventory, payment)

		uuid, err := service.CreateOrder(context.Background(), nil)

		assert.Nil(t, uuid)
		assert.Empty(t, repo.Calls)
		assert.Empty(t, inventory.Calls)
		assert.Empty(t, payment.Calls)
		assert.Equal(t, fmt.Errorf("order info is nil"), err)
	})

	t.Run("Error get parts", func(t *testing.T) {
		repo := repoMocks.NewMockOrderRepository(t)
		inventory := grpcMocks.NewMockIntentoryServiceClient(t)
		payment := grpcMocks.NewMockPaymentServiceClient(t)
		service := New(repo, inventory, payment)

		inventory.EXPECT().ListParts(context.Background(), mock.Anything).Return(nil, fmt.Errorf("error")).Once()
		uuid, err := service.CreateOrder(context.Background(), &model.OrderInfo{
			UserUuid: uuid.NewString(),
			PartIds:  []string{uuid.NewString()},
		})

		assert.Nil(t, uuid)
		assert.Empty(t, repo.Calls)
		assert.Empty(t, payment.Calls)
		assert.Equal(t, fmt.Errorf("error"), err)
	})

	t.Run("Not enough parts", func(t *testing.T) {
		repo := repoMocks.NewMockOrderRepository(t)
		inventory := grpcMocks.NewMockIntentoryServiceClient(t)
		payment := grpcMocks.NewMockPaymentServiceClient(t)
		service := New(repo, inventory, payment)

		parts := make([]*model.Part, 0, 5)
		for range cap(parts) {
			parts = append(parts, &model.Part{})
		}
		inventory.EXPECT().ListParts(context.Background(), mock.Anything).Return(parts, nil).Once()

		orderInfo := &model.OrderInfo{}
		orderInfo.UserUuid = uuid.NewString()
		for range 10 {
			orderInfo.PartIds = append(orderInfo.PartIds, uuid.NewString())
		}
		uuid, err := service.CreateOrder(context.Background(), orderInfo)

		assert.Nil(t, uuid)
		assert.Equal(t, model.ErrPartNotFound, err)

		assert.Empty(t, repo.Calls)
		assert.Empty(t, payment.Calls)
	})

	t.Run("Nil parts", func(t *testing.T) {
		repo := repoMocks.NewMockOrderRepository(t)
		inventory := grpcMocks.NewMockIntentoryServiceClient(t)
		payment := grpcMocks.NewMockPaymentServiceClient(t)
		service := New(repo, inventory, payment)

		inventory.EXPECT().ListParts(context.Background(), mock.Anything).Return(make([]*model.Part, 0), nil).Once()

		orderInfo := &model.OrderInfo{}
		orderInfo.UserUuid = uuid.NewString()
		for range 10 {
			orderInfo.PartIds = append(orderInfo.PartIds, uuid.NewString())
		}
		uuid, err := service.CreateOrder(context.Background(), orderInfo)

		assert.Nil(t, uuid)
		assert.Equal(t, model.ErrPartNotFound, err)

		assert.Empty(t, repo.Calls)
		assert.Empty(t, payment.Calls)
	})

	t.Run("Error create order", func(t *testing.T) {
		repo := repoMocks.NewMockOrderRepository(t)
		inventory := grpcMocks.NewMockIntentoryServiceClient(t)
		payment := grpcMocks.NewMockPaymentServiceClient(t)
		service := New(repo, inventory, payment)

		inventory.EXPECT().ListParts(context.Background(), mock.Anything).Return(nil, nil).Once()
		repo.EXPECT().CreateOrder(context.Background(), mock.Anything).Return(nil, fmt.Errorf("error")).Once()

		uuid, err := service.CreateOrder(context.Background(), &model.OrderInfo{
			UserUuid: uuid.NewString(),
		})

		assert.Nil(t, uuid)
		assert.Equal(t, fmt.Errorf("error"), err)

		assert.Empty(t, payment.Calls)
	})

	t.Run("Success", func(t *testing.T) {
		repo := repoMocks.NewMockOrderRepository(t)
		inventory := grpcMocks.NewMockIntentoryServiceClient(t)
		payment := grpcMocks.NewMockPaymentServiceClient(t)
		service := New(repo, inventory, payment)

		parts := make([]*model.Part, 0, 5)
		for range cap(parts) {
			parts = append(parts, &model.Part{
				Info: model.PartInfo{
					Price: 10,
				},
			})
		}
		inventory.EXPECT().ListParts(context.Background(), mock.Anything).Return(parts, nil).Once()

		info := &model.OrderInfo{
			UserUuid: uuid.NewString(),
		}

		expectUuid := uuid.NewString()
		expectedOrder := &model.Order{
			Uuid: expectUuid,
			Info: *info,
		}
		repo.EXPECT().CreateOrder(context.Background(), info).Return(expectedOrder, nil).Once()

		order, err := service.CreateOrder(context.Background(), info)

		assert.Equal(t, expectedOrder, order)
		assert.NoError(t, err)
		assert.Equal(t, info.TotalPrice, 50.0)

		assert.Empty(t, payment.Calls)
	})
}

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

func TestPayOrder_errorFromPaymentService(t *testing.T) {
	repo := repoMocks.NewMockOrderRepository(t)
	inventory := grpcMocks.NewMockIntentoryServiceClient(t)
	payment := grpcMocks.NewMockPaymentServiceClient(t)
	service := NewService(repo, inventory, payment)

	userUuid := uuid.NewString()
	orderUuid := uuid.NewString()
	paymentMethod := model.OrderPaymentMethod_CARD
	paymentError := fmt.Errorf("error")

	repo.EXPECT().GetOrderStatus(context.Background(), orderUuid).
		Return(model.OrderStatus_PENDING_PAYMENT, nil).Once()

	payment.EXPECT().PayOrder(context.Background(), orderUuid, userUuid, paymentMethod).
		Return(nil, paymentError).Once()

	uuid, err := service.PayOrder(context.Background(), orderUuid, userUuid, paymentMethod)

	assert.Nil(t, uuid)
	assert.Equal(t, err, paymentError)

	assert.Empty(t, inventory.Calls)
}

func TestPayOrder_errorGetStatusFromRepository(t *testing.T) {
	repo := repoMocks.NewMockOrderRepository(t)
	inventory := grpcMocks.NewMockIntentoryServiceClient(t)
	payment := grpcMocks.NewMockPaymentServiceClient(t)
	service := NewService(repo, inventory, payment)

	userUuid := uuid.NewString()
	orderUuid := uuid.NewString()
	paymentMethod := model.OrderPaymentMethod_CARD
	repoError := fmt.Errorf("error")
	// transactionUuid := uuid.NewString()

	repo.EXPECT().GetOrderStatus(context.Background(), orderUuid).
		Return(model.OrderStatus_UNSPECIFIED, repoError).Once()

	uuid, err := service.PayOrder(context.Background(), orderUuid, userUuid, paymentMethod)

	assert.Nil(t, uuid)
	assert.Equal(t, err, repoError)
}

func TestPayOrder_errorPayFromRepository(t *testing.T) {
	repo := repoMocks.NewMockOrderRepository(t)
	inventory := grpcMocks.NewMockIntentoryServiceClient(t)
	payment := grpcMocks.NewMockPaymentServiceClient(t)
	service := NewService(repo, inventory, payment)

	userUuid := uuid.NewString()
	orderUuid := uuid.NewString()
	paymentMethod := model.OrderPaymentMethod_CARD
	repoError := fmt.Errorf("error")
	transactionUuid := uuid.NewString()

	repo.EXPECT().GetOrderStatus(context.Background(), orderUuid).
		Return(model.OrderStatus_PENDING_PAYMENT, nil).Once()

	repo.EXPECT().PayOrder(context.Background(), orderUuid, transactionUuid, paymentMethod).
		Return(repoError).Once()

	payment.EXPECT().PayOrder(context.Background(), orderUuid, userUuid, paymentMethod).
		Return(&transactionUuid, nil).Once()

	uuid, err := service.PayOrder(context.Background(), orderUuid, userUuid, paymentMethod)

	assert.Nil(t, uuid)
	assert.Equal(t, err, repoError)
}

func TestPayOrder_success(t *testing.T) {
	repo := repoMocks.NewMockOrderRepository(t)
	inventory := grpcMocks.NewMockIntentoryServiceClient(t)
	payment := grpcMocks.NewMockPaymentServiceClient(t)
	service := NewService(repo, inventory, payment)

	userUuid := uuid.NewString()
	orderUuid := uuid.NewString()
	paymentMethod := model.OrderPaymentMethod_CARD
	transactionUuid := uuid.NewString()

	repo.EXPECT().GetOrderStatus(context.Background(), orderUuid).
		Return(model.OrderStatus_PENDING_PAYMENT, nil).Once()
	payment.EXPECT().PayOrder(context.Background(), orderUuid, userUuid, paymentMethod).
		Return(&transactionUuid, nil).Once()
	repo.EXPECT().PayOrder(context.Background(), orderUuid, transactionUuid, paymentMethod).
		Return(nil).Once()

	uuid, err := service.PayOrder(context.Background(), orderUuid, userUuid, paymentMethod)

	assert.NoError(t, err)
	assert.Equal(t, uuid, &transactionUuid)

	assert.Empty(t, inventory.Calls)
}

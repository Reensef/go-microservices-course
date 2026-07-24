package order

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	grpcMocks "github.com/Reensef/go-microservices-course/order/internal/client/grpc/mocks"
	eventMocks "github.com/Reensef/go-microservices-course/order/internal/events/mocks"
	"github.com/Reensef/go-microservices-course/order/internal/model"
	repoMocks "github.com/Reensef/go-microservices-course/order/internal/repository/mocks"
)

func TestPayOrder_errorFromPaymentService(t *testing.T) {
	repo := repoMocks.NewMockOrderRepository(t)
	inventory := grpcMocks.NewMockIntentoryClient(t)
	payment := grpcMocks.NewMockPaymentClient(t)
	producer := eventMocks.NewMockOrderProducer(t)
	service := New(repo, inventory, payment, producer)

	userUuid := uuid.NewString()
	orderUuid := uuid.NewString()
	paymentMethod := model.OrderPaymentMethod_CARD
	paymentError := fmt.Errorf("error")

	order := &model.Order{}

	repo.EXPECT().GetOrderByUUID(context.Background(), orderUuid).
		Return(order, nil).Once()

	payment.EXPECT().PayOrder(context.Background(), orderUuid, userUuid, paymentMethod).
		Return("", paymentError).Once()

	uuid, err := service.PayOrder(context.Background(), orderUuid, userUuid, paymentMethod)

	assert.Empty(t, uuid)
	assert.Equal(t, err, paymentError)

	assert.Empty(t, inventory.Calls)
}

func TestPayOrder_errorPayFromRepository(t *testing.T) {
	repo := repoMocks.NewMockOrderRepository(t)
	inventory := grpcMocks.NewMockIntentoryClient(t)
	payment := grpcMocks.NewMockPaymentClient(t)
	producer := eventMocks.NewMockOrderProducer(t)
	service := New(repo, inventory, payment, producer)

	userUuid := uuid.NewString()
	orderUuid := uuid.NewString()
	paymentMethod := model.OrderPaymentMethod_CARD
	repoError := fmt.Errorf("error")
	transactionUuid := uuid.NewString()

	order := &model.Order{}

	repo.EXPECT().GetOrderByUUID(context.Background(), orderUuid).
		Return(order, nil).Once()

	repo.EXPECT().PayOrder(context.Background(), orderUuid, transactionUuid, paymentMethod).
		Return(repoError).Once()

	payment.EXPECT().PayOrder(context.Background(), orderUuid, userUuid, paymentMethod).
		Return(transactionUuid, nil).Once()

	uuid, err := service.PayOrder(context.Background(), orderUuid, userUuid, paymentMethod)

	assert.Empty(t, uuid)
	assert.Equal(t, err, repoError)
}

func TestPayOrder_success(t *testing.T) {
	repo := repoMocks.NewMockOrderRepository(t)
	inventory := grpcMocks.NewMockIntentoryClient(t)
	payment := grpcMocks.NewMockPaymentClient(t)
	producer := eventMocks.NewMockOrderProducer(t)
	service := New(repo, inventory, payment, producer)

	userUuid := uuid.NewString()
	orderUuid := uuid.NewString()
	paymentMethod := model.OrderPaymentMethod_CARD
	transactionUuid := uuid.NewString()

	order := &model.Order{}

	repo.EXPECT().GetOrderByUUID(context.Background(), orderUuid).
		Return(order, nil).Once()

	payment.EXPECT().PayOrder(context.Background(), orderUuid, userUuid, paymentMethod).
		Return(transactionUuid, nil).Once()

	repo.EXPECT().PayOrder(context.Background(), orderUuid, transactionUuid, paymentMethod).
		Return(nil).Once()

	producer.EXPECT().ProduceOrderPaid(context.Background(), mock.Anything).
		Return(nil).Once()

	uuid, err := service.PayOrder(context.Background(), orderUuid, userUuid, paymentMethod)

	assert.NoError(t, err)
	assert.Equal(t, transactionUuid, uuid)

	assert.Empty(t, inventory.Calls)
}

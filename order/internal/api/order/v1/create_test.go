package v1

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Reensef/go-microservices-course/order/internal/model"
	"github.com/Reensef/go-microservices-course/order/internal/service/mocks"
	orderApi "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
)

func TestCreateOrder_usesAuthenticatedUser(t *testing.T) {
	service := mocks.NewMockOrderService(t)
	a := NewHandler(service)

	userUuid := uuid.NewString()
	partId := "part-1"

	ctx := model.WithUser(t.Context(), model.User{Uuid: userUuid})

	service.EXPECT().
		CreateOrder(ctx, &model.OrderInfo{UserUuid: userUuid, PartIds: []string{partId}}).
		Return(&model.Order{Uuid: uuid.NewString(), Info: model.OrderInfo{UserUuid: userUuid}}, nil).
		Once()

	resp, err := a.CreateOrder(ctx, &orderApi.CreateOrderRequest{PartIds: []string{partId}})

	assert.NoError(t, err)
	_, ok := resp.(*orderApi.CreateOrderResponse)
	assert.True(t, ok)
}

func TestCreateOrder_noAuthenticatedUser(t *testing.T) {
	service := mocks.NewMockOrderService(t)
	a := NewHandler(service)

	resp, err := a.CreateOrder(t.Context(), &orderApi.CreateOrderRequest{PartIds: []string{"part-1"}})

	assert.NoError(t, err)
	_, ok := resp.(*orderApi.ForbiddenError)
	assert.True(t, ok)
	assert.Empty(t, service.Calls)
}

func TestCreateOrder_serviceError(t *testing.T) {
	service := mocks.NewMockOrderService(t)
	a := NewHandler(service)

	userUuid := uuid.NewString()

	ctx := model.WithUser(t.Context(), model.User{Uuid: userUuid})

	service.EXPECT().
		CreateOrder(ctx, &model.OrderInfo{UserUuid: userUuid, PartIds: []string{"part-1"}}).
		Return(nil, fmt.Errorf("boom")).
		Once()

	resp, err := a.CreateOrder(ctx, &orderApi.CreateOrderRequest{PartIds: []string{"part-1"}})

	assert.NoError(t, err)
	_, ok := resp.(*orderApi.InternalServerError)
	assert.True(t, ok)
}

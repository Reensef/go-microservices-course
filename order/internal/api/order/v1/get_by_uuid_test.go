package v1

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Reensef/go-microservices-course/order/internal/model"
	"github.com/Reensef/go-microservices-course/order/internal/service/mocks"
	orderApi "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
)

func TestGetOrderByUUID_ownOrder(t *testing.T) {
	service := mocks.NewMockOrderService(t)
	a := NewHandler(service)

	userUuid := uuid.NewString()
	orderUuid := uuid.NewString()

	order := &model.Order{Uuid: orderUuid, Info: model.OrderInfo{UserUuid: userUuid}}

	ctx := model.WithUser(t.Context(), model.User{Uuid: userUuid})

	service.EXPECT().
		GetOrderByUUID(ctx, orderUuid).
		Return(order, nil).
		Once()

	resp, err := a.GetOrderByUUID(ctx, orderApi.GetOrderByUUIDParams{OrderUUID: orderUuid})

	assert.NoError(t, err)
	_, ok := resp.(*orderApi.OrderDto)
	assert.True(t, ok)
}

func TestGetOrderByUUID_noAuthenticatedUser(t *testing.T) {
	service := mocks.NewMockOrderService(t)
	a := NewHandler(service)

	orderUuid := uuid.NewString()

	resp, err := a.GetOrderByUUID(t.Context(), orderApi.GetOrderByUUIDParams{OrderUUID: orderUuid})

	assert.NoError(t, err)
	_, ok := resp.(*orderApi.ForbiddenError)
	assert.True(t, ok)
	assert.Empty(t, service.Calls)
}

func TestGetOrderByUUID_accessDenied(t *testing.T) {
	service := mocks.NewMockOrderService(t)
	a := NewHandler(service)

	userUuid := uuid.NewString()
	ownerUuid := uuid.NewString()
	orderUuid := uuid.NewString()

	order := &model.Order{Uuid: orderUuid, Info: model.OrderInfo{UserUuid: ownerUuid}}

	ctx := model.WithUser(t.Context(), model.User{Uuid: userUuid})

	service.EXPECT().
		GetOrderByUUID(ctx, orderUuid).
		Return(order, nil).
		Once()

	resp, err := a.GetOrderByUUID(ctx, orderApi.GetOrderByUUIDParams{OrderUUID: orderUuid})

	assert.NoError(t, err)
	forbidden, ok := resp.(*orderApi.ForbiddenError)
	assert.True(t, ok)
	assert.Equal(t, 403, forbidden.Code)
}

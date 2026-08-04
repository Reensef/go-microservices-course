package v1

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Reensef/go-microservices-course/order/internal/model"
	"github.com/Reensef/go-microservices-course/order/internal/service/mocks"
	orderApi "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
)

func TestCancelOrder_usesAuthenticatedUser(t *testing.T) {
	service := mocks.NewMockOrderService(t)
	a := NewHandler(service)

	userUuid := uuid.NewString()
	orderUuid := uuid.NewString()

	ctx := model.WithUser(t.Context(), model.User{Uuid: userUuid})

	service.EXPECT().
		CancelOrder(ctx, orderUuid, userUuid).
		Return(nil).
		Once()

	resp, err := a.CancelOrder(ctx, orderApi.CancelOrderParams{OrderUUID: orderUuid})

	assert.NoError(t, err)
	_, ok := resp.(*orderApi.CancelOrderNoContent)
	assert.True(t, ok)
}

func TestCancelOrder_noAuthenticatedUser(t *testing.T) {
	service := mocks.NewMockOrderService(t)
	a := NewHandler(service)

	resp, err := a.CancelOrder(t.Context(), orderApi.CancelOrderParams{OrderUUID: uuid.NewString()})

	assert.NoError(t, err)
	_, ok := resp.(*orderApi.ForbiddenError)
	assert.True(t, ok)
	assert.Empty(t, service.Calls)
}

func TestCancelOrder_accessDenied(t *testing.T) {
	service := mocks.NewMockOrderService(t)
	a := NewHandler(service)

	userUuid := uuid.NewString()
	orderUuid := uuid.NewString()

	ctx := model.WithUser(t.Context(), model.User{Uuid: userUuid})

	service.EXPECT().
		CancelOrder(ctx, orderUuid, userUuid).
		Return(model.ErrOrderAccessDenied).
		Once()

	resp, err := a.CancelOrder(ctx, orderApi.CancelOrderParams{OrderUUID: orderUuid})

	assert.NoError(t, err)
	forbidden, ok := resp.(*orderApi.ForbiddenError)
	assert.True(t, ok)
	assert.Equal(t, 403, forbidden.Code)
}

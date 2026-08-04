package v1

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Reensef/go-microservices-course/order/internal/model"
	"github.com/Reensef/go-microservices-course/order/internal/service/mocks"
	orderApi "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
)

func TestPayOrder_usesAuthenticatedUser(t *testing.T) {
	service := mocks.NewMockOrderService(t)
	a := NewHandler(service)

	userUuid := uuid.NewString()
	orderUuid := uuid.NewString()
	transactionUuid := uuid.NewString()

	ctx := model.WithUser(t.Context(), model.User{Uuid: userUuid})

	service.EXPECT().
		PayOrder(ctx, orderUuid, userUuid, model.OrderPaymentMethod_CARD).
		Return(transactionUuid, nil).
		Once()

	resp, err := a.PayOrder(
		ctx,
		&orderApi.PayOrderRequest{PaymentMethod: orderApi.PaymentMethodCARD},
		orderApi.PayOrderParams{OrderUUID: orderUuid},
	)

	assert.NoError(t, err)
	payResp, ok := resp.(*orderApi.PayOrderResponse)
	assert.True(t, ok)
	assert.Equal(t, transactionUuid, payResp.TransactionUUID.Value)
}

func TestPayOrder_noAuthenticatedUser(t *testing.T) {
	service := mocks.NewMockOrderService(t)
	a := NewHandler(service)

	resp, err := a.PayOrder(
		t.Context(),
		&orderApi.PayOrderRequest{PaymentMethod: orderApi.PaymentMethodCARD},
		orderApi.PayOrderParams{OrderUUID: uuid.NewString()},
	)

	assert.NoError(t, err)
	_, ok := resp.(*orderApi.ForbiddenError)
	assert.True(t, ok)
	assert.Empty(t, service.Calls)
}

func TestPayOrder_accessDenied(t *testing.T) {
	service := mocks.NewMockOrderService(t)
	a := NewHandler(service)

	userUuid := uuid.NewString()
	orderUuid := uuid.NewString()

	ctx := model.WithUser(t.Context(), model.User{Uuid: userUuid})

	service.EXPECT().
		PayOrder(ctx, orderUuid, userUuid, model.OrderPaymentMethod_CARD).
		Return("", model.ErrOrderAccessDenied).
		Once()

	resp, err := a.PayOrder(
		ctx,
		&orderApi.PayOrderRequest{PaymentMethod: orderApi.PaymentMethodCARD},
		orderApi.PayOrderParams{OrderUUID: orderUuid},
	)

	assert.NoError(t, err)
	forbidden, ok := resp.(*orderApi.ForbiddenError)
	assert.True(t, ok)
	assert.Equal(t, 403, forbidden.Code)
}

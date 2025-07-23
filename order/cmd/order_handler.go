package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	orderV1 "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/payment/v1"
)

type OrderHandler struct {
	storage          *OrderStorage
	inventoryService inventoryV1.InventoryServiceClient
	paymentService   paymentV1.PaymentServiceClient
}

// NewOrderHandler создает новый обработчик запросов к API погоды
func NewOrderHandler(storage *OrderStorage,
	inventoryService inventoryV1.InventoryServiceClient,
	paymentService paymentV1.PaymentServiceClient,
) *OrderHandler {
	return &OrderHandler{
		storage:          storage,
		inventoryService: inventoryService,
		paymentService:   paymentService,
	}
}

func (o *OrderHandler) CancelOrder(
	ctx context.Context,
	params orderV1.CancelOrderParams,
) (orderV1.CancelOrderRes, error) {
	order := o.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("Order by UUID '%s' not found", params.OrderUUID.String()),
		}, nil
	}

	switch order.Status {
	case orderV1.OrderStatusPENDINGPAYMENT:
		order.Status = orderV1.OrderStatusCANCELED

		err := o.storage.UpdateOrder(order)
		if err != nil {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Order with UUID '%s' not found", params.OrderUUID.String()),
			}, nil
		}

		return &orderV1.CancelOrderNoContent{}, nil
	case orderV1.OrderStatusPAID:
		return &orderV1.ConflictError{
			Code:    409,
			Message: fmt.Sprintf("Order with UUID '%s' already paid", params.OrderUUID.String()),
		}, nil
	case orderV1.OrderStatusCANCELED:
		return &orderV1.CancelOrderNoContent{}, nil
	default:
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}
}

func (o *OrderHandler) CreateOrder(
	ctx context.Context,
	req *orderV1.CreateOrderRequest,
) (orderV1.CreateOrderRes, error) {
	reqPartUuids := req.GetPartUuids()
	reqPartUuidStrings := make([]string, len(reqPartUuids))
	for i, uuid := range reqPartUuids {
		reqPartUuidStrings[i] = uuid.String()
	}

	resListParts, err := o.inventoryService.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids: reqPartUuidStrings,
		},
	})
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	partUuids := map[string]bool{}
	for _, part := range resListParts.GetParts() {
		partUuids[part.GetUuid()] = true
	}

	for _, uuid := range reqPartUuidStrings {
		if !partUuids[uuid] {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Part with UUID '%s' not found", uuid),
			}, nil
		}
	}

	totalPrice := 0.0
	for _, part := range resListParts.GetParts() {
		totalPrice += part.GetPrice()
	}

	newOrder := &orderV1.OrderDto{
		UserUUID:   req.UserUUID,
		PartUuids:  req.PartUuids,
		TotalPrice: totalPrice,
		Status:     orderV1.OrderStatusPENDINGPAYMENT,
	}
	o.storage.AddOrder(newOrder)

	return &orderV1.CreateOrderResponse{
		OrderUUID:  newOrder.OrderUUID,
		TotalPrice: totalPrice,
	}, nil
}

// Обрабатывает запрос на получение данных о заказе по UUID
func (o *OrderHandler) GetOrderByUUID(_ context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error) {
	order := o.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("Order with UUID '%s' not found", params.OrderUUID.String()),
		}, nil
	}

	return order, nil
}

func (o *OrderHandler) PayOrder(
	ctx context.Context,
	req *orderV1.PayOrderRequest,
	params orderV1.PayOrderParams,
) (orderV1.PayOrderRes, error) {
	order := o.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("Order with UUID '%s' not found", params.OrderUUID.String()),
		}, nil
	}

	orderPaid, err := o.paymentService.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid: order.OrderUUID.String(),
		UserUuid:  order.UserUUID.String(),
	})
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	transactionUUID, err := uuid.Parse(orderPaid.TransactionUuid)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	order.TransactionUUID.SetTo(transactionUUID)
	order.PaymentMethod.SetTo(req.PaymentMethod)
	order.Status = orderV1.OrderStatusPAID

	err = o.storage.UpdateOrder(order)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	response := &orderV1.PayOrderResponse{}
	response.TransactionUUID.SetTo(transactionUUID)

	return response, nil
}

// Cоздает новую ошибку в формате GenericError
func (o *OrderHandler) NewError(
	_ context.Context,
	err error,
) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		},
	}
}

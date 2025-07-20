package main

import (
	"context"
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
			Message: "Order by UUID '" + params.OrderUUID.String() + "' not found",
		}, nil
	}

	switch order.Status {
	case orderV1.OrderStatusPENDINGPAYMENT:
		order.Status = orderV1.OrderStatusCANCELED

		err := o.storage.UpdateOrder(order)
		if err != nil {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: "Order by UUID '" + params.OrderUUID.String() + "' not found",
			}, nil
		}

		return &orderV1.CancelOrderNoContent{}, nil
	case orderV1.OrderStatusPAID:
		return &orderV1.ConflictError{
			Code:    409,
			Message: "Order '" + params.OrderUUID.String() + "' is already paid",
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

	for _, uuid := range reqPartUuidStrings {
		exists := false
		for _, part := range resListParts.GetParts() {
			if part.GetUuid() == uuid {
				exists = true
				break
			}
		}
		if !exists {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: "Part with UUID '" + uuid + "' not found",
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
			Message: "Order by UUID '" + params.OrderUUID.String() + "' not found",
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
			Message: "Order by uuid '" + params.OrderUUID.String() + "' not found",
		}, nil
	}

	resPayOrder, err := o.paymentService.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid: order.OrderUUID.String(),
		UserUuid:  order.UserUUID.String(),
	})
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	transactionUUID, err := uuid.Parse(resPayOrder.TransactionUuid)
	if err != nil {
		return nil, err
	}

	order.TransactionUUID.SetTo(transactionUUID)
	order.PaymentMethod.SetTo(req.PaymentMethod)
	order.Status = orderV1.OrderStatusPAID

	err = o.storage.UpdateOrder(order)
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order by uuid '" + params.OrderUUID.String() + "' not found",
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

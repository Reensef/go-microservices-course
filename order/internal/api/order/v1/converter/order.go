package converter

import (
	"github.com/google/uuid"

	"github.com/Reensef/go-microservices-course/order/internal/model"
	orderV1 "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
)

func ModelOrderToAPI(order *model.Order) *orderV1.OrderDto {
	return &orderV1.OrderDto{
		OrderUUID:       order.Uuid,
		UserUUID:        order.Info.UserUuid,
		PartUuids:       append([]uuid.UUID(nil), order.Info.PartUuids...),
		TotalPrice:      order.Info.TotalPrice,
		TransactionUUID: orderV1.NewOptUUID(order.Info.TransactionUuid),
		PaymentMethod:   ModelPaymentMethodToAPI(order.Info.PaymentMethod),
		Status:          ModelOrderStatusToAPI(order.Info.Status),
	}
}

func APIOrderToModel(order *orderV1.OrderDto) *model.Order {
	return &model.Order{
		Uuid: order.OrderUUID,
		Info: model.OrderInfo{
			UserUuid:      order.UserUUID,
			PartUuids:     order.PartUuids,
			TotalPrice:    order.TotalPrice,
			PaymentMethod: APIPaymentMethodToModel(order.PaymentMethod),
			Status:        APIOrderStatusToModel(order.Status),
		},
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}
}

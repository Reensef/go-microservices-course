package converter

import (
	"github.com/Reensef/go-microservices-course/order/internal/model"
	orderV1 "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
)

func ToApiOrder(order *model.Order) *orderV1.OrderDto {
	return &orderV1.OrderDto{
		OrderUUID:       order.Uuid,
		UserUUID:        order.Info.UserUuid,
		PartIds:         order.Info.PartIds,
		TotalPrice:      order.Info.TotalPrice,
		TransactionUUID: orderV1.NewOptString(order.Info.TransactionUuid),
		PaymentMethod:   ToApiPaymentMethod(order.Info.PaymentMethod),
		Status:          ToApiOrderStatus(order.Info.Status),
	}
}

func ToModelOrder(order *orderV1.OrderDto) *model.Order {
	return &model.Order{
		Uuid: order.OrderUUID,
		Info: model.OrderInfo{
			UserUuid:      order.UserUUID,
			PartIds:       order.PartIds,
			TotalPrice:    order.TotalPrice,
			PaymentMethod: ToModelPaymentMethod(order.PaymentMethod),
			Status:        ToModelOrderStatus(order.Status),
		},
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}
}

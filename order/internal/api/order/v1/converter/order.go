package converter

import (
	"github.com/Reensef/go-microservices-course/order/internal/model"
	orderApi "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
)

func ToApiOrder(order *model.Order) *orderApi.OrderDto {
	return &orderApi.OrderDto{
		OrderUUID:       order.Uuid,
		UserUUID:        order.Info.UserUuid,
		PartIds:         order.Info.PartIds,
		TotalPrice:      order.Info.TotalPrice,
		TransactionUUID: orderApi.NewOptString(order.Info.TransactionUuid),
		PaymentMethod:   ToApiPaymentMethod(order.Info.PaymentMethod),
		Status:          ToApiOrderStatus(order.Info.Status),
	}
}

func ToModelOrder(order *orderApi.OrderDto) *model.Order {
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

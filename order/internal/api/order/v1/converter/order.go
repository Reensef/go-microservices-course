package converter

import (
	"github.com/Reensef/go-microservices-course/order/internal/model"
	orderV1 "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
	"github.com/gofrs/uuid"
)

func ModelOrderToAPI(order model.Order) orderV1.OrderDto {
	return orderV1.OrderDto{
		OrderUUID:       order.Uuid,
		UserUUID:        order.Info.UserUuid,
		PartUuids:       append([]uuid.UUID(nil), order.Info.PartUuids...),
		TotalPrice:      order.Info.TotalPrice,
		TransactionUUID: orderV1.NewOptUUID(order.Info.TransactionUuid),
		PaymentMethod:   order.Info.PaymentMethod,
		Status:          order.Info.Status,
	}
}

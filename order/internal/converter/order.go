package converter

import (
	"github.com/Reensef/go-microservices-course/order/internal/model"
	repoModel "github.com/Reensef/go-microservices-course/order/internal/repository/model"
		orderV1 "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"

	"github.com/google/uuid"
)

func ModelOrderToAPI(order model.Order) orderV1.OrderDto {
	return orderV1.OrderDto{
		OrderUUID: order.Uuid,
		UserUUID: order.Info.UserUuid,
		PartUuids: append([]uuid.UUID(nil), order.Info.PartUuids...),
		TotalPrice: order.Info.TotalPrice,
		TransactionUUID: orderV1.NewOptUUID(order.Info.TransactionUuid),
		PaymentMethod: order.Info.PaymentMethod,
		Status: order.Info.Status,

	}
}

func ModelOrderToRepoModel(order model.Order) repoModel.Order {
	return repoModel.Order{
		Uuid:      order.Uuid,
		Info:      ModelOrderInfoToRepoModel(order.Info),
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
		DeletedAt: order.UpdatedAt,
	}
}

func RepoModelOrderInfoToModel(info repoModel.OrderInfo) model.OrderInfo {
	return model.OrderInfo{
		UserUuid:        info.UserUuid,
		PartUuids:       append([]uuid.UUID(nil), info.PartUuids...),
		TransactionUuid: info.TransactionUuid,
		TotalPrice:      info.TotalPrice,
		PaymentMethod:   info.PaymentMethod,
		Status:          info.Status,
	}
}

func RepoModelOrderToModel(order repoModel.Order) model.Order {
	return model.Order{
		Uuid:      order.Uuid,
		Info:      RepoModelOrderInfoToModel(order.Info),
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
		DeletedAt: order.DeletedAt,
	}
}

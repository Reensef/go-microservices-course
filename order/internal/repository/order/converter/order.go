package converter

import (
	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/Reensef/go-microservices-course/order/internal/model"
	repoModel "github.com/Reensef/go-microservices-course/order/internal/repository/model"
)

func ModelOrderInfoToRepoModel(info *model.OrderInfo) *repoModel.OrderInfo {
	return &repoModel.OrderInfo{
		UserUuid:        info.UserUuid,
		PartUuids:       append([]uuid.UUID(nil), info.PartUuids...),
		TransactionUuid: info.TransactionUuid,
		TotalPrice:      info.TotalPrice,
		PaymentMethod:   ModelPaymentMethodToRepo(info.PaymentMethod),
		Status:          ModelOrderStatusToRepo(info.Status),
	}
}

func ModelOrderToRepoModel(order *model.Order) *repoModel.Order {
	return &repoModel.Order{
		Uuid:      order.Uuid,
		Info:      *ModelOrderInfoToRepoModel(&order.Info),
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
		DeletedAt: lo.ToPtr(*order.DeletedAt),
	}
}

func RepoModelOrderInfoToModel(info *repoModel.OrderInfo) *model.OrderInfo {
	return &model.OrderInfo{
		UserUuid:        info.UserUuid,
		PartUuids:       append([]uuid.UUID(nil), info.PartUuids...),
		TransactionUuid: info.TransactionUuid,
		TotalPrice:      info.TotalPrice,
		PaymentMethod:   RepoPaymentMethodToModel(info.PaymentMethod),
		Status:          RepoOrderStatusToModel(info.Status),
	}
}

func RepoModelOrderToModel(order *repoModel.Order) *model.Order {
	return &model.Order{
		Uuid:      order.Uuid,
		Info:      *RepoModelOrderInfoToModel(&order.Info),
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
		DeletedAt: order.DeletedAt,
	}
}

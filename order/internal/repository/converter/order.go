package converter

import (
	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/Reensef/go-microservices-course/order/internal/model"
	repoModel "github.com/Reensef/go-microservices-course/order/internal/repository/model"
)

func ToRepoModelOrder(order *model.Order) *repoModel.Order {
	if order == nil {
		return nil
	}
	return &repoModel.Order{
		Uuid:      order.Uuid,
		Info:      *ToRepoModelOrderInfo(&order.Info),
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
		DeletedAt: lo.ToPtr(*order.DeletedAt),
	}
}

func ToRepoModelOrderInfo(info *model.OrderInfo) *repoModel.OrderInfo {
	if info == nil {
		return nil
	}
	return &repoModel.OrderInfo{
		UserUuid:        info.UserUuid,
		PartUuids:       append([]uuid.UUID(nil), info.PartUuids...),
		TransactionUuid: info.TransactionUuid,
		TotalPrice:      info.TotalPrice,
		PaymentMethod:   ToRepoModelPaymentMethod(info.PaymentMethod),
		Status:          ToRepoModelOrderStatus(info.Status),
	}
}

func ToModelOrder(order *repoModel.Order) *model.Order {
	if order == nil {
		return nil
	}
	return &model.Order{
		Uuid:      order.Uuid,
		Info:      *ToModelOrderInfo(&order.Info),
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
		DeletedAt: order.DeletedAt,
	}
}

func ToModelOrderInfo(info *repoModel.OrderInfo) *model.OrderInfo {
	if info == nil {
		return nil
	}
	return &model.OrderInfo{
		UserUuid:        info.UserUuid,
		PartUuids:       append([]uuid.UUID(nil), info.PartUuids...),
		TransactionUuid: info.TransactionUuid,
		TotalPrice:      info.TotalPrice,
		PaymentMethod:   ToModelPaymentMethod(info.PaymentMethod),
		Status:          ToModelOrderStatus(info.Status),
	}
}

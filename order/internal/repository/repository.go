package repository

import (
	"context"

	repoModel "github.com/Reensef/go-microservices-course/order/internal/repository/model"
	"github.com/google/uuid"
)

type OrderRepository interface {
	CreateOrder(
		ctx context.Context,
		userUuid uuid.UUID,
		partUuids []uuid.UUID,
	) (repoModel.Order, error)

	GetOrderByUUID(
		ctx context.Context,
		uuid uuid.UUID,
	) (repoModel.Order, error)

	CancelOrder(
		ctx context.Context,
		uuid uuid.UUID,
	) error

	PayOrder(
		ctx context.Context,
		uuid uuid.UUID,
		paymentMethod repoModel.OrderPaymentMethod,
	) (uuid.UUID, error)
}

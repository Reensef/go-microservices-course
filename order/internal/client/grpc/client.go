package grpc

import (
	"context"

	"github.com/google/uuid"

	"github.com/Reensef/go-microservices-course/order/internal/model"
)

type IntentoryServiceClient interface {
	GetPart(ctx context.Context, partUuid uuid.UUID) (*model.Part, error)
	ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error)
}

type PaymentServiceClient interface {
	PayOrder(
		ctx context.Context,
		orderUuid, userUuid uuid.UUID,
		paymentMethod model.OrderPaymentMethod,
	) (transactionUuid *uuid.UUID, err error)
}

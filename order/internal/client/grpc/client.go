package grpc

import (
	"context"

	"github.com/Reensef/go-microservices-course/order/internal/model"
)

type IntentoryClient interface {
	GetPart(ctx context.Context, partUuid string) (*model.Part, error)
	ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error)
}

type PaymentClient interface {
	PayOrder(
		ctx context.Context,
		orderUuid, userUuid string,
		paymentMethod model.OrderPaymentMethod,
	) (transactionUuid string, err error)
}

type IAMClient interface {
	Whoami(ctx context.Context, sessionUUID string) (model.User, error)
}

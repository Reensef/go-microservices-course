package order

import (
	"context"

	"github.com/google/uuid"

	model "github.com/Reensef/go-microservices-course/order/internal/model"
	repoConverter "github.com/Reensef/go-microservices-course/order/internal/repository/converter"
	repoModel "github.com/Reensef/go-microservices-course/order/internal/repository/model"
)

func (r *repository) PayOrder(
	ctx context.Context,
	orderUuid uuid.UUID,
	transactionUUID uuid.UUID,
	paymentMethod model.OrderPaymentMethod,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.data[orderUuid]
	if !ok {
		return model.ErrOrderNotFound
	}

	order.Info.PaymentMethod = repoConverter.ToRepoModelPaymentMethod(paymentMethod)
	order.Info.Status = repoModel.OrderStatus_PAID
	order.Info.TransactionUuid = transactionUUID

	return nil
}

package order

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	model "github.com/Reensef/go-microservices-course/order/internal/model"
	repoConverter "github.com/Reensef/go-microservices-course/order/internal/repository/converter"
	repoModel "github.com/Reensef/go-microservices-course/order/internal/repository/model"
)

func (r *repository) PayOrder(
	ctx context.Context,
	orderUuid string,
	transactionUUID string,
	paymentMethod model.OrderPaymentMethod,
) error {
	builderUpdate := sq.Update("orders").
		Set("order_status", repoModel.OrderStatus_PAID).
		Set("payment_method", repoConverter.ToRepoModelPaymentMethod(paymentMethod)).
		Set("transaction_uuid", transactionUUID).
		Where(sq.Eq{"uuid": orderUuid}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return err
	}

	res, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return model.ErrOrderNotFound
	}

	return nil
}

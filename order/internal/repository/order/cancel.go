package order

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	model "github.com/Reensef/go-microservices-course/order/internal/model"
	repoModel "github.com/Reensef/go-microservices-course/order/internal/repository/model"
)

func (r *repository) CancelOrder(
	ctx context.Context,
	orderUuid string,
) error {
	builderUpdate := sq.Update("orders").
		Set("order_status", repoModel.OrderStatus_CANCELED).
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

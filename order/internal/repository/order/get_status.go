package order

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"

	model "github.com/Reensef/go-microservices-course/order/internal/model"
	"github.com/Reensef/go-microservices-course/order/internal/repository/converter"
	repoModel "github.com/Reensef/go-microservices-course/order/internal/repository/model"
)

func (r *repository) GetOrderStatus(
	ctx context.Context,
	orderUuid string,
) (model.OrderStatus, error) {
	builderSelect := sq.Select("order_status").
		From("orders").
		Where(sq.Eq{"uuid": orderUuid}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return model.OrderStatus_UNSPECIFIED, err
	}

	row := r.pool.QueryRow(ctx, query, args...)

	var status repoModel.OrderStatus
	err = row.Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.OrderStatus_UNSPECIFIED, errors.Join(model.ErrOrderNotFound, err)
		}
		return model.OrderStatus_UNSPECIFIED, err
	}

	return converter.ToModelOrderStatus(status), nil
}

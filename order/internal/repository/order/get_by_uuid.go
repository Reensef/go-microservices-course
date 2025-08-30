package order

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	model "github.com/Reensef/go-microservices-course/order/internal/model"
	repoConverter "github.com/Reensef/go-microservices-course/order/internal/repository/converter"
	repoModel "github.com/Reensef/go-microservices-course/order/internal/repository/model"
)

func (r *repository) GetOrderByUUID(
	ctx context.Context,
	orderUuid string,
) (*model.Order, error) {
	builderSelect := sq.Select(
		"user_uuid", "part_ids", "transaction_uuid",
		"total_price", "payment_method", "order_status").
		PlaceholderFormat(sq.Dollar).
		From("orders").
		Where(sq.Eq{"uuid": orderUuid})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, err
	}

	order := &repoModel.Order{
		Uuid: orderUuid,
	}

	row := r.pool.QueryRow(ctx, query, args...)
	err = row.Scan(
		&order.Info.UserUuid,
		&order.Info.PartUuids,
		&order.Info.TransactionUuid,
		&order.Info.TotalPrice,
		&order.Info.PaymentMethod,
		&order.Info.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrOrderNotFound
		}

		return nil, err
	}

	return repoConverter.ToModelOrder(order), nil
}

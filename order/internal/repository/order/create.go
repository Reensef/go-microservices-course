package order

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"

	model "github.com/Reensef/go-microservices-course/order/internal/model"
)

func (r *repository) CreateOrder(
	ctx context.Context,
	info *model.OrderInfo,
) (*model.Order, error) {
	partIds := make([]string, 0)
	if info.PartIds != nil {
		partIds = info.PartIds
	}
	builderInsert := sq.Insert("orders").
		PlaceholderFormat(sq.Dollar).
		Columns("user_uuid", "part_ids", "total_price", "payment_method", "order_status").
		Values(info.UserUuid, partIds, info.TotalPrice, info.PaymentMethod, info.Status).
		Suffix("RETURNING uuid, created_at, updated_at")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return nil, err
	}

	var (
		uuid      string
		createdAt time.Time
		updatedAt time.Time
	)

	err = r.pool.QueryRow(ctx, query, args...).Scan(&uuid, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	return &model.Order{
		Uuid:      uuid,
		Info:      *info,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

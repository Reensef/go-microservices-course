package user

import (
	"context"
	"encoding/json"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	model "github.com/Reensef/go-microservices-course/iam/internal/model"
	repoConverter "github.com/Reensef/go-microservices-course/iam/internal/repository/converter"
	repoModel "github.com/Reensef/go-microservices-course/iam/internal/repository/model"
)

func (r *repository) GetByUUID(
	ctx context.Context,
	userUuid string,
) (model.User, error) {
	builderSelect := sq.Select("login", "email", "notification_methods", "created_at", "updated_at").
		PlaceholderFormat(sq.Dollar).
		From("users").
		Where(sq.Eq{"uuid": userUuid})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return model.User{}, err
	}

	var notificationMethods []byte
	user := repoModel.User{Uuid: userUuid}

	row := r.pool.QueryRow(ctx, query, args...)
	err = row.Scan(&user.Info.Login, &user.Info.Email, &notificationMethods, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, model.ErrUserNotFound
		}

		return model.User{}, err
	}

	if err := json.Unmarshal(notificationMethods, &user.Info.NotificationMethods); err != nil {
		return model.User{}, err
	}

	return repoConverter.ToModelUser(user), nil
}

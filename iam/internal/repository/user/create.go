package user

import (
	"context"
	"encoding/json"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"

	model "github.com/Reensef/go-microservices-course/iam/internal/model"
	repoConverter "github.com/Reensef/go-microservices-course/iam/internal/repository/converter"
	repoModel "github.com/Reensef/go-microservices-course/iam/internal/repository/model"
)

const uniqueViolationCode = "23505"

func (r *repository) Create(
	ctx context.Context,
	info *model.UserRegistrationInfo,
	passwordHash string,
) (model.User, error) {
	notificationMethods, err := json.Marshal(repoConverter.ToRepoNotificationMethods(info.Info.NotificationMethods))
	if err != nil {
		return model.User{}, err
	}

	builderInsert := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("login", "email", "password_hash", "notification_methods").
		Values(info.Info.Login, info.Info.Email, passwordHash, notificationMethods).
		Suffix("RETURNING uuid, created_at, updated_at")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return model.User{}, err
	}

	user := repoModel.User{
		Info: repoModel.UserInfo{
			Login:               info.Info.Login,
			Email:               info.Info.Email,
			NotificationMethods: repoConverter.ToRepoNotificationMethods(info.Info.NotificationMethods),
		},
	}

	err = r.pool.QueryRow(ctx, query, args...).Scan(&user.Uuid, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return model.User{}, model.ErrUserAlreadyExists
		}

		return model.User{}, err
	}

	return repoConverter.ToModelUser(user), nil
}

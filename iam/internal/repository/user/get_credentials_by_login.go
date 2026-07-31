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

func (r *repository) GetCredentialsByLogin(
	ctx context.Context,
	login string,
) (*model.User, string, error) {
	builderSelect := sq.Select("uuid", "email", "password_hash", "notification_methods", "created_at", "updated_at").
		PlaceholderFormat(sq.Dollar).
		From("users").
		Where(sq.Eq{"login": login})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, "", err
	}

	var (
		notificationMethods []byte
		passwordHash        string
	)
	user := repoModel.User{Info: repoModel.UserInfo{Login: login}}

	row := r.pool.QueryRow(ctx, query, args...)
	err = row.Scan(&user.Uuid, &user.Info.Email, &passwordHash, &notificationMethods, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, "", model.ErrUserNotFound
		}

		return nil, "", err
	}

	if err := json.Unmarshal(notificationMethods, &user.Info.NotificationMethods); err != nil {
		return nil, "", err
	}

	result := repoConverter.ToModelUser(user)
	return &result, passwordHash, nil
}

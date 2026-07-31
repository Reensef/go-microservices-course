package session

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/redis/go-redis/v9"

	model "github.com/Reensef/go-microservices-course/iam/internal/model"
	repoConverter "github.com/Reensef/go-microservices-course/iam/internal/repository/converter"
	repoModel "github.com/Reensef/go-microservices-course/iam/internal/repository/model"
)

func (r *repository) GetByUUID(
	ctx context.Context,
	sessionUuid string,
) (*model.Session, error) {
	data, err := r.client.Get(ctx, key(sessionUuid)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, model.ErrSessionNotFound
		}

		return nil, err
	}

	var session repoModel.Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}

	result := repoConverter.ToModelSession(session)
	return &result, nil
}

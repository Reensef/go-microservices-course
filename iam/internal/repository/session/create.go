package session

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	model "github.com/Reensef/go-microservices-course/iam/internal/model"
	repoConverter "github.com/Reensef/go-microservices-course/iam/internal/repository/converter"
	repoModel "github.com/Reensef/go-microservices-course/iam/internal/repository/model"
)

func (r *repository) Create(
	ctx context.Context,
	userUuid string,
	ttl time.Duration,
) (model.Session, error) {
	now := time.Now().UTC()

	session := repoModel.Session{
		Uuid:      uuid.NewString(),
		UserUuid:  userUuid,
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: now.Add(ttl),
	}

	data, err := json.Marshal(session)
	if err != nil {
		return model.Session{}, err
	}

	if err := r.client.Set(ctx, key(session.Uuid), data, ttl).Err(); err != nil {
		return model.Session{}, err
	}

	return repoConverter.ToModelSession(session), nil
}

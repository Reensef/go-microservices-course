package session

import (
	"github.com/redis/go-redis/v9"

	repo "github.com/Reensef/go-microservices-course/iam/internal/repository"
)

const keyPrefix = "iam:session:"

var _ repo.SessionRepository = (*repository)(nil)

type repository struct {
	client *redis.Client
}

func New(client *redis.Client) *repository {
	return &repository{
		client: client,
	}
}

func key(sessionUuid string) string {
	return keyPrefix + sessionUuid
}

package user

import (
	"github.com/jackc/pgx/v5/pgxpool"

	repo "github.com/Reensef/go-microservices-course/iam/internal/repository"
)

var _ repo.UserRepository = (*repository)(nil)

type repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *repository {
	return &repository{
		pool: pool,
	}
}

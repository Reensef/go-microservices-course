package order

import (
	"sync"

	"github.com/google/uuid"

	repo "github.com/Reensef/go-microservices-course/order/internal/repository"
	repoModel "github.com/Reensef/go-microservices-course/order/internal/repository/model"
)

var _ repo.OrderRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data map[uuid.UUID]*repoModel.Order
}

func NewRepository() *repository {
	return &repository{
		data: make(map[uuid.UUID]*repoModel.Order),
	}
}

package grpc

import (
	"context"

	"github.com/google/uuid"

	"github.com/Reensef/go-microservices-course/order/internal/model"
)

type PartsFilter struct {
	Uuids                 []string
	Names                 []string
	Categories            []model.PartCategory
	ManufacturerCountries []string
	Tags                  []string
}

type IntentoryServiceClient interface {
	// Метод для получения детали по UUID
	GetPart(ctx context.Context, partUuid *uuid.UUID) (*model.Part, error)
	// Метод для получения списка деталей с фильтрацией
	ListParts(ctx context.Context, filter *PartsFilter) ([]*model.Part, error)
}

type PaymentServiceClient interface {
	// Метод для оплаты заказа
	// Возвращает UUID транзакции
	PayOrder(ctx context.Context, orderUuid, userUuid *uuid.UUID) (*uuid.UUID, error)
}

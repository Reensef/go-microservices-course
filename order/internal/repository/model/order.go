package model

import (
	"time"

	"github.com/Reensef/go-microservices-course/order/internal/model"

	"github.com/google/uuid"
)

type OrderInfo struct {
	UserUuid        uuid.UUID
	PartUuids       []uuid.UUID
	TransactionUuid uuid.UUID
	TotalPrice      float64
	PaymentMethod   model.OrderPaymentMethod
	Status          model.OrderStatus
}

type OrderUpdateInfo struct {
	UserUuid        *uuid.UUID
	PartUuids       []uuid.UUID
	TransactionUuid *uuid.UUID
	TotalPrice      *float64
	PaymentMethod   *model.OrderPaymentMethod
	Status          *model.OrderStatus
}

type Order struct {
	Uuid      uuid.UUID
	Info      OrderInfo
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

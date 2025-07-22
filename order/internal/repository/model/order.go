package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderPaymentMethod int32

const (
	OrderPaymentMethod_UNSPECIFIED    OrderPaymentMethod = 0
	OrderPaymentMethod_CARD           OrderPaymentMethod = 1
	OrderPaymentMethod_CREDIT_CARD    OrderPaymentMethod = 2
	OrderPaymentMethod_SBP            OrderPaymentMethod = 3
	OrderPaymentMethod_INVESTOR_MONEY OrderPaymentMethod = 4
)

type OrderStatus int

const (
	OrderStatus_PENDING_PAYMENT OrderStatus = 0
	OrderStatus_PAID            OrderStatus = 1
	OrderStatus_CANCELED        OrderStatus = 2
)

type OrderInfo struct {
	UserUuid        uuid.UUID
	PartUuids       []uuid.UUID
	TransactionUuid uuid.UUID
	TotalPrice      float64
	PaymentMethod   OrderPaymentMethod
	Status          OrderStatus
}

type OrderUpdateInfo struct {
	UserUuid        *uuid.UUID
	PartUuids       []uuid.UUID
	TransactionUuid *uuid.UUID
	TotalPrice      *float64
	PaymentMethod   *OrderPaymentMethod
	Status          *OrderStatus
}

type Order struct {
	Uuid      uuid.UUID
	Info      OrderInfo
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

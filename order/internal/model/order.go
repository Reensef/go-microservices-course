package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderPaymentMethod int

const (
	OrderPaymentMethod_UNSPECIFIED = iota
	OrderPaymentMethod_CARD
	OrderPaymentMethod_CREDIT_CARD
	OrderPaymentMethod_SBP
	OrderPaymentMethod_INVESTOR_MONEY
)

type OrderStatus int

const (
	OrderStatus_PENDING_PAYMENT = iota
	OrderStatus_PAID
	OrderStatus_CANCELED
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

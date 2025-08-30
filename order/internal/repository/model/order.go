package model

import (
	"database/sql"
	"time"
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
	OrderStatus_UNSPECIFIED     OrderStatus = 0
	OrderStatus_PENDING_PAYMENT OrderStatus = 1
	OrderStatus_PAID            OrderStatus = 2
	OrderStatus_CANCELED        OrderStatus = 3
)

type OrderInfo struct {
	UserUuid        string
	PartUuids       []string
	TransactionUuid sql.NullString
	TotalPrice      float64
	PaymentMethod   OrderPaymentMethod
	Status          OrderStatus
}

type OrderUpdateInfo struct {
	UserUuid        *string
	PartUuids       []string
	TransactionUuid *string
	TotalPrice      *float64
	PaymentMethod   *OrderPaymentMethod
	Status          *OrderStatus
}

type Order struct {
	Uuid      string
	Info      OrderInfo
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

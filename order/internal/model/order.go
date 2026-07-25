package model

import (
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

func (m OrderPaymentMethod) String() string {
	switch m {
	case OrderPaymentMethod_UNSPECIFIED:
		return "UNSPECIFIED"
	case OrderPaymentMethod_CARD:
		return "CARD"
	case OrderPaymentMethod_CREDIT_CARD:
		return "CREDIT_CARD"
	case OrderPaymentMethod_SBP:
		return "SBP"
	case OrderPaymentMethod_INVESTOR_MONEY:
		return "INVESTOR_MONEY"
	default:
		return "UNKNOWN"
	}
}

type OrderStatus int32

const (
	OrderStatus_UNSPECIFIED     OrderStatus = 0
	OrderStatus_PENDING_PAYMENT OrderStatus = 1
	OrderStatus_PAID            OrderStatus = 2
	OrderStatus_CANCELED        OrderStatus = 3
	OrderStatus_ASSEMBLED       OrderStatus = 4
)

type OrderInfo struct {
	UserUuid        string
	PartIds         []string
	TransactionUuid string
	TotalPrice      float64
	PaymentMethod   OrderPaymentMethod
	Status          OrderStatus
}

type OrderUpdateInfo struct {
	UserUuid        *string
	PartIds         []string
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

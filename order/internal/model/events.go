package model

import "time"

type ShipAssembledEvent struct {
	UUID          string
	OrderUUID     string
	UserUUID      string
	BuildDuration time.Duration
}

type OrderPaidEvent struct {
	UUID            string
	OrderUUID       string
	UserUUID        string
	TransactionUUID string
	PaymentMethod   OrderPaymentMethod
}

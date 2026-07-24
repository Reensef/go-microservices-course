package model

import "time"

type OrderPaidEvent struct {
	UUID            string
	OrderUUID       string
	UserUUID        string
	TransactionUUID string
	PaymentMethod   string
}

type ShipAssembledEvent struct {
	UUID          string
	OrderUUID     string
	UserUUID      string
	BuildDuration time.Duration
}

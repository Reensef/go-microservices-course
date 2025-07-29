package model

type PaymentMethod int32

const (
	PaymentMethod_UNSPECIFIED PaymentMethod = iota
	PaymentMethod_CARD
	PaymentMethod_SBP
	PaymentMethod_CREDIT_CARD
	PaymentMethod_INVESTOR_MONEY
)

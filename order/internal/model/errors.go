package model

import "errors"

var (
	ErrOrderNotFound    = errors.New("order not found")
	ErrPartNotFound     = errors.New("part not found")
	ErrOrderAlreadyPaid = errors.New("order already paid")

	ErrOrderUuidInvalidFormat = errors.New("order UUID must be UUID format")
	ErrUserUuidInvalidFormat  = errors.New("user UUID must be UUID format")
	ErrPartIdInvalidFormat    = errors.New("part ID must be ObjectID format")

	ErrPaymentMethodUnspecified = errors.New("payment method unspecified")
)

package model

import "errors"

var (
	ErrPaymentMethodUnspecified = errors.New("payment method unspecified")
	ErrOrderUuidInvalidFormat   = errors.New("order UUID must be UUID format")
	ErrUserUuidInvalidFormat    = errors.New("user UUID must be UUID format")
)

package model

import "errors"

var ErrOrderNotFound = errors.New("order not found")
var ErrPartNotFound = errors.New("part not found")
var ErrOrderAlreadyPaid = errors.New("order already paid")

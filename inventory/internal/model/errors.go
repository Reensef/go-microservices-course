package model

import "errors"

var (
	ErrPartNotFound        = errors.New("part not found")
	ErrPartIdInvalidFormat = errors.New("part id must be ObjectID format")
)

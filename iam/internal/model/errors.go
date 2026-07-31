package model

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user with this login already exists")
	ErrInvalidCredentials = errors.New("invalid login or password")

	ErrSessionNotFound = errors.New("session not found")

	ErrUserUuidInvalidFormat    = errors.New("user UUID must be UUID format")
	ErrSessionUuidInvalidFormat = errors.New("session UUID must be UUID format")
)

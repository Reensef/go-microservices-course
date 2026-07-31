package model

import "time"

type Session struct {
	Uuid      string
	UserUuid  string
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}

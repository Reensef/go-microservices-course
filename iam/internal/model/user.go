package model

import "time"

type NotificationMethod struct {
	ProviderName string
	Target       string
}

type UserInfo struct {
	Login               string
	Email               string
	NotificationMethods []NotificationMethod
}

type UserRegistrationInfo struct {
	Info     UserInfo
	Password string
}

type User struct {
	Uuid      string
	Info      UserInfo
	CreatedAt time.Time
	UpdatedAt time.Time
}

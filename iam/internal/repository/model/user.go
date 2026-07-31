package model

import "time"

type NotificationMethod struct {
	ProviderName string `json:"provider_name"`
	Target       string `json:"target"`
}

type UserInfo struct {
	Login               string
	Email               string
	NotificationMethods []NotificationMethod
}

type User struct {
	Uuid      string
	Info      UserInfo
	CreatedAt time.Time
	UpdatedAt time.Time
}

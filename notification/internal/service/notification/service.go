package notification

import (
	client "github.com/Reensef/go-microservices-course/notification/internal/client/telegram"
	def "github.com/Reensef/go-microservices-course/notification/internal/service"
)

var _ def.NotificationService = (*service)(nil)

type service struct {
	telegramClient client.Client
	chatID         int64
}

func New(telegramClient client.Client, chatID int64) *service {
	return &service{
		telegramClient: telegramClient,
		chatID:         chatID,
	}
}

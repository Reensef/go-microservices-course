package notification

import (
	"context"

	"github.com/Reensef/go-microservices-course/notification/internal/model"
)

func (s *service) NotifyShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error {
	text, err := renderTemplate(shipAssembledTemplate, event)
	if err != nil {
		return err
	}

	return s.telegramClient.SendMessage(ctx, s.chatID, text)
}

package notification

import (
	"context"

	"github.com/Reensef/go-microservices-course/notification/internal/model"
)

func (s *service) NotifyOrderPaid(ctx context.Context, event model.OrderPaidEvent) error {
	text, err := renderTemplate(orderPaidTemplate, event)
	if err != nil {
		return err
	}

	return s.telegramClient.SendMessage(ctx, s.chatID, text)
}

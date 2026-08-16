package v1

import (
	"context"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"github.com/Reensef/go-microservices-course/platform/pkg/logger"
)

const startMessage = "👋 Привет! Я бот уведомлений интернет-магазина.\n" +
	"Я пришлю сообщение сюда, когда заказ будет оплачен или собран."

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// RegisterHandlers регистрирует обработчики команд бота.
func (h *Handler) RegisterHandlers(b *tgbot.Bot) {
	b.RegisterHandler(tgbot.HandlerTypeMessageText, "start", tgbot.MatchTypeCommand, h.HandleStart)
}

// HandleStart отвечает на команду /start приветственным сообщением.
func (h *Handler) HandleStart(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	_, err := b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   startMessage,
	})
	if err != nil {
		logger.Error("failed to send /start reply", zap.Error(err))
	}
}

package bot

import (
	"context"

	tgbot "github.com/go-telegram/bot"

	def "github.com/Reensef/go-microservices-course/notification/internal/client/telegram"
)

var _ def.Client = (*client)(nil)

type client struct {
	bot *tgbot.Bot
}

func New(bot *tgbot.Bot) *client {
	return &client{
		bot: bot,
	}
}

func (c *client) SendMessage(ctx context.Context, chatID int64, text string) error {
	_, err := c.bot.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})

	return err
}

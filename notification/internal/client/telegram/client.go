package telegram

import "context"

type Client interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
}

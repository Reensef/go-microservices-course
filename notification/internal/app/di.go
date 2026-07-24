package app

import (
	"context"
	"fmt"

	tgbot "github.com/go-telegram/bot"

	telegramApi "github.com/Reensef/go-microservices-course/notification/internal/api/telegram/v1"
	telegramClient "github.com/Reensef/go-microservices-course/notification/internal/client/telegram"
	telegramBotClient "github.com/Reensef/go-microservices-course/notification/internal/client/telegram/bot"
	"github.com/Reensef/go-microservices-course/notification/internal/config"
	events "github.com/Reensef/go-microservices-course/notification/internal/events"
	notificationConsumer "github.com/Reensef/go-microservices-course/notification/internal/events/kafka/consumer/notification"
	service "github.com/Reensef/go-microservices-course/notification/internal/service"
	notificationService "github.com/Reensef/go-microservices-course/notification/internal/service/notification"
	closer "github.com/Reensef/go-microservices-course/platform/pkg/closer"
)

type diContainer struct {
	notificationConsumer events.Consumer
	notificationService  service.NotificationService

	telegramHandler *telegramApi.Handler
	telegramClient  telegramClient.Client
	telegramBot     *tgbot.Bot
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) NotificationConsumer(ctx context.Context) events.Consumer {
	if d.notificationConsumer == nil {
		consumer, err := notificationConsumer.NewConsumer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().NotificationConsumer.GroupID(),
			config.AppConfig().NotificationConsumer.OrderPaidTopic(),
			config.AppConfig().NotificationConsumer.ShipAssembledTopic(),
			config.AppConfig().NotificationConsumer.Config(),
			d.NotificationService(ctx),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create notification consumer: %s\n", err.Error()))
		}

		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return consumer.Close()
		})

		d.notificationConsumer = consumer
	}

	return d.notificationConsumer
}

func (d *diContainer) NotificationService(ctx context.Context) service.NotificationService {
	if d.notificationService == nil {
		d.notificationService = notificationService.New(
			d.TelegramClient(ctx),
			config.AppConfig().Telegram.ChatID(),
		)
	}

	return d.notificationService
}

func (d *diContainer) TelegramClient(ctx context.Context) telegramClient.Client {
	if d.telegramClient == nil {
		d.telegramClient = telegramBotClient.New(d.TelegramBot(ctx))
	}

	return d.telegramClient
}

func (d *diContainer) TelegramHandler(_ context.Context) *telegramApi.Handler {
	if d.telegramHandler == nil {
		d.telegramHandler = telegramApi.NewHandler()
	}

	return d.telegramHandler
}

func (d *diContainer) TelegramBot(ctx context.Context) *tgbot.Bot {
	if d.telegramBot == nil {
		b, err := tgbot.New(config.AppConfig().Telegram.BotToken())
		if err != nil {
			panic(fmt.Sprintf("failed to create telegram bot: %s\n", err.Error()))
		}

		d.TelegramHandler(ctx).RegisterHandlers(b)

		d.telegramBot = b
	}

	return d.telegramBot
}

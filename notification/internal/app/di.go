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

func (d *diContainer) NotificationConsumer(ctx context.Context) (events.Consumer, error) {
	if d.notificationConsumer == nil {
		svc, err := d.NotificationService(ctx)
		if err != nil {
			return nil, err
		}

		consumer, err := notificationConsumer.NewConsumer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().NotificationConsumer.GroupID(),
			config.AppConfig().NotificationConsumer.OrderPaidTopic(),
			config.AppConfig().NotificationConsumer.ShipAssembledTopic(),
			config.AppConfig().NotificationConsumer.Config(),
			svc,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create notification consumer: %w", err)
		}

		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return consumer.Close()
		})

		d.notificationConsumer = consumer
	}

	return d.notificationConsumer, nil
}

func (d *diContainer) NotificationService(ctx context.Context) (service.NotificationService, error) {
	if d.notificationService == nil {
		client, err := d.TelegramClient(ctx)
		if err != nil {
			return nil, err
		}

		d.notificationService = notificationService.New(
			client,
			config.AppConfig().Telegram.ChatID(),
		)
	}

	return d.notificationService, nil
}

func (d *diContainer) TelegramClient(ctx context.Context) (telegramClient.Client, error) {
	if d.telegramClient == nil {
		bot, err := d.TelegramBot(ctx)
		if err != nil {
			return nil, err
		}

		d.telegramClient = telegramBotClient.New(bot)
	}

	return d.telegramClient, nil
}

func (d *diContainer) TelegramHandler(_ context.Context) *telegramApi.Handler {
	if d.telegramHandler == nil {
		d.telegramHandler = telegramApi.NewHandler()
	}

	return d.telegramHandler
}

func (d *diContainer) TelegramBot(ctx context.Context) (*tgbot.Bot, error) {
	if d.telegramBot == nil {
		b, err := tgbot.New(config.AppConfig().Telegram.BotToken())
		if err != nil {
			return nil, fmt.Errorf("failed to create telegram bot: %w", err)
		}

		d.TelegramHandler(ctx).RegisterHandlers(b)

		d.telegramBot = b
	}

	return d.telegramBot, nil
}

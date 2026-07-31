//go:build integration

package environment

import (
	"context"

	"go.uber.org/zap"

	"github.com/Reensef/go-microservices-course/platform/pkg/logger"
)

// Teardown — освобождает все ресурсы тестового окружения
func Teardown(ctx context.Context, env *TestEnvironment) {
	log := logger.Logger()
	log.Info(ctx, "🧹 Очистка тестового окружения...")

	cleanup(ctx, env)

	log.Info(ctx, "✅ Тестовое окружение успешно очищено")
}

// cleanup — вспомогательная функция для освобождения ресурсов
func cleanup(ctx context.Context, env *TestEnvironment) {
	if env.App != nil {
		if err := env.App.Terminate(ctx); err != nil {
			logger.Error(ctx, "не удалось остановить контейнер приложения", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 Контейнер приложения остановлен")
		}
	}

	if env.Redis != nil {
		if err := env.Redis.Terminate(ctx); err != nil {
			logger.Error(ctx, "не удалось остановить контейнер Redis", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 Контейнер Redis остановлен")
		}
	}

	if env.Postgres != nil {
		if err := env.Postgres.Terminate(ctx); err != nil {
			logger.Error(ctx, "не удалось остановить контейнер Postgres", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 Контейнер Postgres остановлен")
		}
	}

	if env.Network != nil {
		if err := env.Network.Remove(ctx); err != nil {
			logger.Error(ctx, "не удалось удалить сеть", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 Сеть удалена")
		}
	}
}

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
	log.Info("🧹 Очистка тестового окружения...")

	cleanup(ctx, env)

	log.Info("✅ Тестовое окружение успешно очищено")
}

// cleanup — вспомогательная функция для освобождения ресурсов
func cleanup(ctx context.Context, env *TestEnvironment) {
	if env.App != nil {
		if err := env.App.Terminate(ctx); err != nil {
			logger.Error("не удалось остановить контейнер приложения", zap.Error(err))
		} else {
			logger.Info("🛑 Контейнер приложения остановлен")
		}
	}

	if env.Redis != nil {
		if err := env.Redis.Terminate(ctx); err != nil {
			logger.Error("не удалось остановить контейнер Redis", zap.Error(err))
		} else {
			logger.Info("🛑 Контейнер Redis остановлен")
		}
	}

	if env.Postgres != nil {
		if err := env.Postgres.Terminate(ctx); err != nil {
			logger.Error("не удалось остановить контейнер Postgres", zap.Error(err))
		} else {
			logger.Info("🛑 Контейнер Postgres остановлен")
		}
	}

	if env.Network != nil {
		if err := env.Network.Remove(ctx); err != nil {
			logger.Error("не удалось удалить сеть", zap.Error(err))
		} else {
			logger.Info("🛑 Сеть удалена")
		}
	}
}

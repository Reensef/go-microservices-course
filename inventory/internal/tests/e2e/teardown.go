//go:build integration

package integration

import (
	"context"

	"go.uber.org/zap"

	"github.com/Reensef/go-microservices-course/platform/pkg/logger"
)

// teardownTestEnvironment — освобождает все ресурсы тестового окружения
func teardownTestEnvironment(ctx context.Context, env *TestEnvironment) {
	log := logger.Logger()
	log.Info("🧹 Очистка тестового окружения...")

	cleanupTestEnvironment(ctx, env)

	log.Info("✅ Тестовое окружение успешно очищено")
}

// cleanupTestEnvironment — вспомогательная функция для освобождения ресурсов
func cleanupTestEnvironment(ctx context.Context, env *TestEnvironment) {
	if env.App != nil {
		if err := env.App.Terminate(ctx); err != nil {
			logger.Error("не удалось остановить контейнер приложения", zap.Error(err))
		} else {
			logger.Info("🛑 Контейнер приложения остановлен")
		}
	}

	if env.Mongo != nil {
		if err := env.Mongo.Terminate(ctx); err != nil {
			logger.Error("не удалось остановить контейнер MongoDB", zap.Error(err))
		} else {
			logger.Info("🛑 Контейнер MongoDB остановлен")
		}
	}

	if env.IamApp != nil {
		if err := env.IamApp.Terminate(ctx); err != nil {
			logger.Error("не удалось остановить контейнер приложения iam", zap.Error(err))
		} else {
			logger.Info("🛑 Контейнер приложения iam остановлен")
		}
	}

	if env.IamRedis != nil {
		if err := env.IamRedis.Terminate(ctx); err != nil {
			logger.Error("не удалось остановить контейнер Redis для iam", zap.Error(err))
		} else {
			logger.Info("🛑 Контейнер Redis для iam остановлен")
		}
	}

	if env.IamPostgres != nil {
		if err := env.IamPostgres.Terminate(ctx); err != nil {
			logger.Error("не удалось остановить контейнер Postgres для iam", zap.Error(err))
		} else {
			logger.Info("🛑 Контейнер Postgres для iam остановлен")
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

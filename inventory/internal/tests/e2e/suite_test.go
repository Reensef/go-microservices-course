//go:build integration

package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"go.uber.org/zap/zapcore"

	"github.com/Reensef/go-microservices-course/platform/pkg/logger"
)

// Стенд собирает два Docker-образа (inventory-app и iam-app, нужен для auth-interceptor),
// поэтому бюджет вдвое больше, чем у наборов с одним образом (например, iam).
const testsTimeout = 10 * time.Minute

var (
	env *TestEnvironment

	suiteCtx    context.Context
	suiteCancel context.CancelFunc
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UFO Service Integration Test Suite")
}

var _ = BeforeSuite(func() {
	var level zapcore.Level
	_ = level.UnmarshalText([]byte(loggerLevelValue))
	err := logger.Init(
		level,
		logger.WithJSON(true),
	)
	if err != nil {
		panic(fmt.Sprintf("не удалось инициализировать логгер: %v", err))
	}

	suiteCtx, suiteCancel = context.WithTimeout(context.Background(), testsTimeout)

	logger.Info("Запуск тестового окружения...")
	env = setupTestEnvironment(suiteCtx)
})

var _ = AfterSuite(func() {
	logger.Info("Завершение набора тестов")
	if env != nil {
		teardownTestEnvironment(suiteCtx, env)
	}
	suiteCancel()
})

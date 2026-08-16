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

	"github.com/Reensef/go-microservices-course/iam/internal/tests/e2e/environment"
	"github.com/Reensef/go-microservices-course/platform/pkg/logger"
)

const testsTimeout = 5 * time.Minute

var (
	env *environment.TestEnvironment

	suiteCtx    context.Context
	suiteCancel context.CancelFunc
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IAM Service Integration Test Suite")
}

var _ = BeforeSuite(func() {
	err := logger.Init(
		zapcore.DebugLevel,
		logger.WithJSON(true),
	)
	if err != nil {
		panic(fmt.Sprintf("не удалось инициализировать логгер: %v", err))
	}

	suiteCtx, suiteCancel = context.WithTimeout(context.Background(), testsTimeout)

	logger.Info("Запуск тестового окружения...")
	env = environment.Setup(suiteCtx)
})

var _ = AfterSuite(func() {
	logger.Info("Завершение набора тестов")
	if env != nil {
		environment.Teardown(suiteCtx, env)
	}
	suiteCancel()
})

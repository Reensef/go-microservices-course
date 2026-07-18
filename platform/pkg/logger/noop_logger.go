package logger

import (
	"context"

	"go.uber.org/zap"
)

type DummyLogger struct{}

func (l *DummyLogger) Info(ctx context.Context, msg string, fields ...zap.Field)  {}
func (l *DummyLogger) Error(ctx context.Context, msg string, fields ...zap.Field) {}

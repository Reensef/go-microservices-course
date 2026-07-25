package logger

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	core := zapcore.NewNopCore()

	globalLogger = &logger{
		zapLogger: zap.New(core),
	}
}

func BenchmarkGlobalLogger(b *testing.B) {
	ctx := context.Background()

	for b.Loop() {
		Info(ctx, "test message")
	}
}

func BenchmarkWithLogger(b *testing.B) {
	log := With(zap.String("static_field", "static_value"))
	ctx := context.Background()

	for b.Loop() {
		log.Info(ctx, "test message")
	}
}

func BenchmarkWithContextLogger(b *testing.B) {
	ctx := context.WithValue(context.Background(), traceIDKey, "trace-123")
	ctx = context.WithValue(ctx, userIDKey, "user-456")

	for b.Loop() {
		WithContext(ctx).Info(ctx, "test message")
	}
}

func BenchmarkChainLogger(b *testing.B) {
	ctx := context.WithValue(context.Background(), traceIDKey, "trace-123")
	ctx = context.WithValue(ctx, userIDKey, "user-456")

	log := With(zap.String("static_field", "static_value"))

	for b.Loop() {
		log.Info(ctx, "test message")
	}
}

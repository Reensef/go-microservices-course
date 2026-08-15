package logger

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	globalLogger = &logger{zapLogger: zap.New(zapcore.NewNopCore())}
}

func BenchmarkGlobalLogger(b *testing.B) {
	for b.Loop() {
		Info("test message")
	}
}

func BenchmarkWithLogger(b *testing.B) {
	log := With(zap.String("static_field", "static_value"))

	for b.Loop() {
		log.Info("test message")
	}
}

func BenchmarkWithFieldLogger(b *testing.B) {
	log := With(
		zap.String("static_field", "static_value"),
		zap.String("service", "order"),
	)

	for b.Loop() {
		log.Info("test message", zap.String("dynamic", "value"))
	}
}

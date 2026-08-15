package logger

import (
	"context"
	"os"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	otelLog "go.opentelemetry.io/otel/log"
	otelLogSdk "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)


var (
	globalLogger *logger
	initOnce     sync.Once
	dynamicLevel zap.AtomicLevel
	otelProvider *otelLogSdk.LoggerProvider
)

type logger struct {
	zapLogger *zap.Logger
}

// init устанавливает no-op логгер по умолчанию — вызовы до Init не паникуют.
func init() {
	globalLogger = &logger{zapLogger: zap.NewNop()}
}

// Option задаёт параметр инициализации логгера.
type Option func(*initConfig)

type initConfig struct {
	level           zapcore.Level
	asJSON          bool
	otlpEndpoint    string
	otlpServiceName string
	otlpEnvironment string
}

// WithLevel задаёт уровень логирования. По умолчанию "info".
func WithLevel(level zapcore.Level) Option {
	return func(c *initConfig) { c.level = level }
}

// WithJSON включает JSON-формат вывода. По умолчанию консольный формат.
func WithJSON(enabled bool) Option {
	return func(c *initConfig) { c.asJSON = enabled }
}

// WithOTLP включает экспорт логов в OpenTelemetry коллектор через gRPC.
// Если подключение не удаётся, логгер продолжает работать только со stdout.
func WithOTLP(endpoint, serviceName, environment string) Option {
	return func(c *initConfig) {
		c.otlpEndpoint = endpoint
		c.otlpServiceName = serviceName
		c.otlpEnvironment = environment
	}
}

// Init инициализирует глобальный логгер
func Init(level zapcore.Level, opts ...Option) error {
	cfg := &initConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	cfg.level = level

	initOnce.Do(func() {
		stdoutCore := newStdoutCore(cfg)
		cores := []zapcore.Core{stdoutCore}

		var otlpErr error
		if cfg.otlpEndpoint != "" {
			otlpCore, err := newOTLPCore(cfg)
			if err != nil {
				otlpErr = err
			} else {
				cores = append(cores, otlpCore)
			}
		}

		var core zapcore.Core
		if len(cores) == 1 {
			core = cores[0]
		} else {
			core = zapcore.NewTee(cores...)
		}

		globalLogger = &logger{
			zapLogger: zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2)),
		}

		if otlpErr != nil {
			globalLogger.zapLogger.Warn("OTLP init failed, continuing without it", zap.Error(otlpErr))
		}
	})
	return nil
}

func newStdoutCore(cfg *initConfig) zapcore.Core {
	dynamicLevel = zap.NewAtomicLevelAt(cfg.level)

	encCfg := defaultEncoderConfig()
	var encoder zapcore.Encoder
	if cfg.asJSON {
		encoder = zapcore.NewJSONEncoder(encCfg)
	} else {
		encoder = zapcore.NewConsoleEncoder(encCfg)
	}

	return zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), dynamicLevel)
}

func newOTLPCore(cfg *initConfig) (*simpleOTLPCore, error) {
	ctx := context.Background()

	exporter, err := otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(cfg.otlpEndpoint),
		otlploggrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	rs, err := resource.New(ctx,
		resource.WithAttributes(
			attribute.String("service.name", cfg.otlpServiceName),
			attribute.String("deployment.environment", cfg.otlpEnvironment),
		),
	)
	if err != nil {
		return nil, err
	}

	provider := otelLogSdk.NewLoggerProvider(
		otelLogSdk.WithResource(rs),
		otelLogSdk.WithProcessor(otelLogSdk.NewBatchProcessor(exporter)),
	)
	otelProvider = provider

	var otlpLogger otelLog.Logger = provider.Logger("app")
	return newSimpleOTLPCore(otlpLogger, dynamicLevel), nil
}

// SetLevel динамически меняет уровень логирования.
func SetLevel(level zapcore.Level) {
	if dynamicLevel == (zap.AtomicLevel{}) {
		return
	}
	dynamicLevel.SetLevel(level)
}

// Logger возвращает глобальный логгер.
func Logger() *logger {
	return globalLogger
}

// Sync сбрасывает буферы логгера.
func Sync() error {
	if globalLogger != nil {
		return globalLogger.zapLogger.Sync()
	}
	return nil
}

// Close корректно завершает работу OTLP provider. Вызывать при graceful shutdown.
func Close(ctx context.Context) error {
	if otelProvider != nil {
		return otelProvider.Shutdown(ctx)
	}
	return nil
}

// defaultEncoderConfig возвращает стандартный EncoderConfig платформы.
func defaultEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
}

// With создаёт логгер с дополнительными постоянными полями.
func With(fields ...zap.Field) *logger {
	if globalLogger == nil {
		return &logger{zapLogger: zap.NewNop()}
	}
	return &logger{zapLogger: globalLogger.zapLogger.With(fields...)}
}

func Debug(msg string, fields ...zap.Field) {
	globalLogger.zapLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	globalLogger.zapLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	globalLogger.zapLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	globalLogger.zapLogger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	globalLogger.zapLogger.Fatal(msg, fields...)
}

func (l *logger) Debug(msg string, fields ...zap.Field) {
	l.zapLogger.Debug(msg, fields...)
}

func (l *logger) Info(msg string, fields ...zap.Field) {
	l.zapLogger.Info(msg, fields...)
}

func (l *logger) Warn(msg string, fields ...zap.Field) {
	l.zapLogger.Warn(msg, fields...)
}

func (l *logger) Error(msg string, fields ...zap.Field) {
	l.zapLogger.Error(msg, fields...)
}

func (l *logger) Fatal(msg string, fields ...zap.Field) {
	l.zapLogger.Fatal(msg, fields...)
}

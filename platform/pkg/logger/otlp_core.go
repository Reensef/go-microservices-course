package logger

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	otelLog "go.opentelemetry.io/otel/log"
	"go.uber.org/zap/zapcore"
)

const emitTimeout = 500 * time.Millisecond

// SimpleOTLPCore конвертирует zap-записи в OpenTelemetry Records и отправляет их в OTLP коллектор.
type SimpleOTLPCore struct {
	otlpLogger otelLog.Logger
	level      zapcore.LevelEnabler
	fields     []zapcore.Field // аккумулированные поля из With()
}

// NewSimpleOTLPCore создает новый OTLP core.
func NewSimpleOTLPCore(otlpLogger otelLog.Logger, level zapcore.LevelEnabler) *SimpleOTLPCore {
	return &SimpleOTLPCore{
		otlpLogger: otlpLogger,
		level:      level,
	}
}

// Enabled проверяет уровень логирования.
func (c *SimpleOTLPCore) Enabled(level zapcore.Level) bool {
	return c.level.Enabled(level)
}

// With создает копию core с дополнительными полями.
// Поля аккумулируются и включаются в каждую запись — в отличие от примера, где они терялись.
func (c *SimpleOTLPCore) With(fields []zapcore.Field) zapcore.Core {
	accumulated := make([]zapcore.Field, len(c.fields)+len(fields))
	copy(accumulated, c.fields)
	copy(accumulated[len(c.fields):], fields)
	return &SimpleOTLPCore{
		otlpLogger: c.otlpLogger,
		level:      c.level,
		fields:     accumulated,
	}
}

// Check добавляет core в список получателей, если уровень подходит.
func (c *SimpleOTLPCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return ce.AddCore(entry, c)
	}
	return ce
}

// Write конвертирует zap Entry + накопленные поля в OTLP Record и отправляет с таймаутом.
func (c *SimpleOTLPCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	severity := mapZapToOtelSeverity(entry.Level)
	record := makeBaseRecord(entry, severity)

	allFields := append(c.fields, fields...)
	if len(allFields) > 0 {
		if attrs := encodeFieldsToAttrs(allFields); len(attrs) > 0 {
			record.AddAttributes(attrs...)
		}
	}

	c.emitWithTimeout(record)
	return nil
}

// Sync — батчинг делает OTLP SDK, явная синхронизация не нужна.
func (c *SimpleOTLPCore) Sync() error { return nil }

func mapZapToOtelSeverity(level zapcore.Level) otelLog.Severity {
	switch level {
	case zapcore.DebugLevel:
		return otelLog.SeverityDebug
	case zapcore.InfoLevel:
		return otelLog.SeverityInfo
	case zapcore.WarnLevel:
		return otelLog.SeverityWarn
	case zapcore.ErrorLevel:
		return otelLog.SeverityError
	default:
		return otelLog.SeverityInfo
	}
}

func makeBaseRecord(entry zapcore.Entry, sev otelLog.Severity) otelLog.Record {
	r := otelLog.Record{}
	r.SetSeverity(sev)
	r.SetBody(attribute.StringValue(entry.Message))
	r.SetTimestamp(entry.Time)
	return r
}

// encodeFieldsToAttrs конвертирует zap поля в OTLP атрибуты.
// Неподдерживаемые типы пропускаются — они продолжат жить в stdout через zap encoder.
func encodeFieldsToAttrs(fields []zapcore.Field) []attribute.KeyValue {
	if len(fields) == 0 {
		return nil
	}

	enc := zapcore.NewMapObjectEncoder()
	for _, f := range fields {
		f.AddTo(enc)
	}

	attrs := make([]attribute.KeyValue, 0, len(enc.Fields))
	for k, v := range enc.Fields {
		switch val := v.(type) {
		case string:
			attrs = append(attrs, attribute.String(k, val))
		case bool:
			attrs = append(attrs, attribute.Bool(k, val))
		case int64:
			attrs = append(attrs, attribute.Int64(k, val))
		case float64:
			attrs = append(attrs, attribute.Float64(k, val))
		}
	}

	return attrs
}

func (c *SimpleOTLPCore) emitWithTimeout(record otelLog.Record) {
	if c.otlpLogger == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), emitTimeout)
	defer cancel()
	c.otlpLogger.Emit(ctx, record)
}

package tracing

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
	"go.opentelemetry.io/otel/trace"
)

type Config interface {
	CollectorEndpoint() string
	ServiceName() string
	Environment() string
	ServiceVersion() string
}

var provider *sdktrace.TracerProvider

func Init(ctx context.Context, cfg Config) error {
	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(cfg.CollectorEndpoint()),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithTimeout(5*time.Second),
		otlptracegrpc.WithCompressor("gzip"),
		otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{
			Enabled:         true,
			InitialInterval: 500 * time.Millisecond,
			MaxInterval:     5 * time.Second,
			MaxElapsedTime:  30 * time.Second,
		}),
	)
	if err != nil {
		return err
	}

	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName()),
			semconv.ServiceVersion(cfg.ServiceVersion()),
			attribute.String("environment", cfg.Environment()),
		),
		resource.WithHost(),
		resource.WithOS(),
		resource.WithProcess(),
		resource.WithContainer(),
		resource.WithTelemetrySDK(),
	)
	if err != nil {
		return err
	}

	provider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.AlwaysSample())),
	)

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return nil
}

func Shutdown(ctx context.Context) error {
	if provider == nil {
		return nil
	}

	return provider.Shutdown(ctx)
}

func StartSpan(
	ctx context.Context, name string, opts ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	return otel.Tracer("").Start(ctx, name, opts...)
}

func RecordError(span trace.Span, err error) {
	if err == nil {
		return
	}

	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

func TraceIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return ""
	}

	return span.SpanContext().TraceID().String()
}

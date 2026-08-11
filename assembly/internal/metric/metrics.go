package metric

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	namespace   = "assembly"
	serviceName = "assembly-service"
)

type metrics struct {
	ordersReceivedCounter     metric.Int64Counter
	ordersProcessedCounter    metric.Int64Counter
	assemblyDurationHistogram metric.Float64Histogram
}

var global *metrics

func Init(ctx context.Context, endpoint string) (*sdkmetric.MeterProvider, error) {
	exporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(endpoint),
		otlpmetricgrpc.WithTLSCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				exporter,
				sdkmetric.WithInterval(5*time.Second),
			),
		),
	)

	otel.SetMeterProvider(meterProvider)

	if err = initInstruments(); err != nil {
		return nil, err
	}

	return meterProvider, nil
}

func initInstruments() error {
	m := otel.Meter(serviceName)

	ordersReceived, err := m.Int64Counter(
		namespace+"_orders_received_total",
		metric.WithDescription("Total number of orders received for assembly"),
	)
	if err != nil {
		return err
	}

	ordersProcessed, err := m.Int64Counter(
		namespace+"_orders_processed_total",
		metric.WithDescription("Total number of orders processed (assembled or failed)"),
	)
	if err != nil {
		return err
	}

	assemblyDuration, err := m.Float64Histogram(
		namespace+"_assembly_duration_seconds",
		metric.WithDescription("Duration of ship assembly"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(
			0.5, 1, 2, 3, 4, 5, 6, 7, 8, 10, 15, 20, 30,
		),
	)
	if err != nil {
		return err
	}

	global = &metrics{
		ordersReceivedCounter:     ordersReceived,
		ordersProcessedCounter:    ordersProcessed,
		assemblyDurationHistogram: assemblyDuration,
	}

	return nil
}

func IncOrdersReceived(ctx context.Context) {
	global.ordersReceivedCounter.Add(ctx, 1)
}

func IncOrdersProcessed(ctx context.Context, status string) {
	global.ordersProcessedCounter.Add(ctx, 1,
		metric.WithAttributes(attribute.String("status", status)),
	)
}

func ObserveAssemblyDuration(ctx context.Context, status string, d float64) {
	global.assemblyDurationHistogram.Record(ctx, d,
		metric.WithAttributes(attribute.String("status", status)),
	)
}

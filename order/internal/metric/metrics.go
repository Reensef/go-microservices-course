package metric

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc/credentials/insecure"
)

type metrics struct {
	ordersCounter        metric.Int64Counter
	ordersRevenueCounter metric.Float64Counter
}

var global *metrics

func init() {
	// no-op instruments until Init is called — prevents nil panics in tests.
	_ = initInstruments("")
}

func Init(ctx context.Context, endpoint, serviceName string) (*sdkmetric.MeterProvider, error) {
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
		resource.WithAttributes(semconv.ServiceName(serviceName)),
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

	err = initInstruments(serviceName)
	if err != nil {
		return nil, err
	}

	return meterProvider, nil
}

func initInstruments(name string) error {
	m := otel.Meter(name)

	ordersCounter, err := m.Int64Counter(
		"orders_total",
		metric.WithDescription("Total number of created orders"),
	)
	if err != nil {
		return err
	}

	ordersRevenueCounter, err := m.Float64Counter(
		"orders_revenue_total",
		metric.WithDescription("Total revenue from created orders"),
		metric.WithUnit("{currency}"),
	)
	if err != nil {
		return err
	}

	global = &metrics{
		ordersCounter:        ordersCounter,
		ordersRevenueCounter: ordersRevenueCounter,
	}

	return nil
}

func IncOrdersTotal(ctx context.Context) {
	global.ordersCounter.Add(ctx, 1)
}

func AddOrderRevenue(ctx context.Context, amount float64) {
	global.ordersRevenueCounter.Add(ctx, amount)
}

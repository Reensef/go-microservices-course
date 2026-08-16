package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const TraceIDHeader = "x-trace-id"

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		propagator := otel.GetTextMapPropagator()
		tracer := otel.Tracer("")

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		ctx = propagator.Extract(ctx, metadataCarrier(md))

		ctx, span := tracer.Start(
			ctx,
			info.FullMethod,
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		setTraceIDHeader(ctx)

		resp, err := handler(ctx, req)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		return resp, err
	}
}

func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		propagator := otel.GetTextMapPropagator()
		tracer := otel.Tracer("")

		ctx, span := tracer.Start(
			ctx,
			method,
			trace.WithSpanKind(trace.SpanKindClient),
		)
		defer span.End()

		carrier := metadataCarrier(extractOutgoingMetadata(ctx))
		propagator.Inject(ctx, carrier)
		ctx = metadata.NewOutgoingContext(ctx, metadata.MD(carrier))

		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		return err
	}
}

// setTraceIDHeader отправляет trace ID клиенту в заголовке ответа gRPC.
func setTraceIDHeader(ctx context.Context) {
	traceID := TraceIDFromContext(ctx)
	if traceID == "" {
		return
	}

	_ = grpc.SetHeader(ctx, metadata.Pairs(TraceIDHeader, traceID))
}

func extractOutgoingMetadata(ctx context.Context) metadata.MD {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return metadata.New(nil)
	}

	return md.Copy()
}

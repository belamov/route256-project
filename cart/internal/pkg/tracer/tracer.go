package tracer

import (
	"context"
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

func NewTraceProvider(exp tracesdk.SpanExporter, ServiceName string) (*tracesdk.TracerProvider, error) {
	// Ensure default SDK resources and the required service name are set.
	r, err := resource.Merge(
		resource.Environment(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(ServiceName),
		),
	)
	if err != nil {
		return nil, err
	}

	return tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(r),
	), nil
}

func InitTracer(ctx context.Context, wg *sync.WaitGroup, endpoint, urlPath string, serviceName string) (trace.Tracer, error) {
	exporter, err := NewSpanExporter(ctx, endpoint, urlPath)
	if err != nil {
		return nil, fmt.Errorf("initialize exporter: %w", err)
	}

	tp, err := NewTraceProvider(exporter, serviceName)
	if err != nil {
		return nil, fmt.Errorf("initialize provider: %w", err)
	}

	go func() {
		defer wg.Done()
		<-ctx.Done()
		log.Info().Msg("shutting down tracer...")
		err = tp.ForceFlush(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("failed to flush tracer provider")
		}
		err = tp.Shutdown(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("failed to shutdown tracer provider")
			return
		}
		log.Info().Msg("tracer stopped")
	}()

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp.Tracer("cart"), nil
}

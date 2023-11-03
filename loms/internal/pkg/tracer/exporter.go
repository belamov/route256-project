package tracer

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func NewSpanExporter(ctx context.Context, endpoint, urlPath string) (tracesdk.SpanExporter, error) {
	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithURLPath(urlPath),
		otlptracehttp.WithInsecure(),
	)
	return otlptrace.New(ctx, client)
}

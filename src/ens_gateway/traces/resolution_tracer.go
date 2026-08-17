package traces

import (
	"context"
	"strconv"

	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/traces/interfaces"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const tracerName = "agentic_ens_gateway"

type ResolutionTracer struct {
	tracer trace.Tracer
}

func NewResolutionTracer() *ResolutionTracer {
	return &ResolutionTracer{tracer: otel.Tracer(tracerName)}
}

func (t *ResolutionTracer) StartResolution(
	ctx context.Context,
	queryKind string,
	chainId uint64,
) (context.Context, interfaces.ResolutionSpan) {
	ctx, span := t.tracer.Start(ctx, "ens.resolve "+queryKind)
	span.SetAttributes(
		attribute.String("ens.query_kind", queryKind),
		attribute.String("chain.id", strconv.FormatUint(chainId, 10)),
	)
	return ctx, &resolutionSpan{otelSpan{span: span}}
}

func (t *ResolutionTracer) StartRegistryRead(ctx context.Context, chainId uint64) (context.Context, interfaces.Span) {
	return t.startChainSpan(ctx, "ens.registry.read", chainId)
}

func (t *ResolutionTracer) StartUniversalResolverFallback(
	ctx context.Context,
	chainId uint64,
) (context.Context, interfaces.Span) {
	return t.startChainSpan(ctx, "ens.universal_resolver.resolve", chainId)
}

func (t *ResolutionTracer) startChainSpan(
	ctx context.Context,
	name string,
	chainId uint64,
) (context.Context, interfaces.Span) {
	ctx, span := t.tracer.Start(ctx, name, trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(attribute.String("chain.id", strconv.FormatUint(chainId, 10)))
	return ctx, &otelSpan{span: span}
}

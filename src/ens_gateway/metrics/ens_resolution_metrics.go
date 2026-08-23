package metrics

import (
	"context"
	"strconv"
	"time"

	common_metrics "github.com/6022-labs/agentic-ens-gateway/src/common/metrics"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

const meterName = "agentic_ens_gateway"

// errorTypeOther is the semconv catch-all for an error with no lower-cardinality classification.
const errorTypeOther = "_OTHER"

type EnsResolutionMetrics struct {
	resolutionHistogram metric.Float64Histogram
	cacheResultCounter  metric.Int64Counter
	fallbackCounter     metric.Int64Counter
	decodeErrorCounter  metric.Int64Counter
	signedCounter       metric.Int64Counter
}

func NewEnsResolutionMetrics() *EnsResolutionMetrics {
	noopMeter := noop.NewMeterProvider().Meter(meterName)
	meter := otel.GetMeterProvider().Meter(meterName)

	resolutionHistogram, err := meter.Float64Histogram(
		"ens.resolution.duration",
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(common_metrics.LatencySecondsBuckets...),
		metric.WithDescription("Duration of CCIP-Read resolutions"),
	)
	if err != nil {
		resolutionHistogram, _ = noopMeter.Float64Histogram("ens.resolution.duration")
	}

	cacheResultCounter, err := meter.Int64Counter(
		"cache.result",
		metric.WithDescription("Resolution cache lookups by hit/miss"),
	)
	if err != nil {
		cacheResultCounter, _ = noopMeter.Int64Counter("cache.result")
	}

	fallbackCounter, err := meter.Int64Counter(
		"ens.fallback",
		metric.WithDescription("L1 universal resolver fallback attempts by outcome"),
	)
	if err != nil {
		fallbackCounter, _ = noopMeter.Int64Counter("ens.fallback")
	}

	decodeErrorCounter, err := meter.Int64Counter(
		"ens.calldata.decode.errors",
		metric.WithDescription("Resolver calldata that could not be decoded"),
	)
	if err != nil {
		decodeErrorCounter, _ = noopMeter.Int64Counter("ens.calldata.decode.errors")
	}

	signedCounter, err := meter.Int64Counter(
		"ens.response.signed",
		metric.WithDescription("Signed CCIP-Read responses by outcome"),
	)
	if err != nil {
		signedCounter, _ = noopMeter.Int64Counter("ens.response.signed")
	}

	return &EnsResolutionMetrics{
		resolutionHistogram: resolutionHistogram,
		cacheResultCounter:  cacheResultCounter,
		fallbackCounter:     fallbackCounter,
		decodeErrorCounter:  decodeErrorCounter,
		signedCounter:       signedCounter,
	}
}

func (m *EnsResolutionMetrics) RecordResolution(
	ctx context.Context,
	queryKind, source string,
	chainId uint64,
	duration time.Duration,
	err error,
) {
	attrs := []attribute.KeyValue{
		attribute.String("ens.query_kind", queryKind),
		attribute.String("ens.source", source),
		attribute.String("chain.id", strconv.FormatUint(chainId, 10)),
	}
	// error.type is conditionally required, and only on a failed operation.
	if err != nil {
		attrs = append(attrs, attribute.String("error.type", errorTypeOther))
	}

	m.resolutionHistogram.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
}

func (m *EnsResolutionMetrics) RecordCacheResult(ctx context.Context, queryKind string, hit bool) {
	result := "miss"
	if hit {
		result = "hit"
	}
	m.cacheResultCounter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("cache.name", "ens_resolution"),
		attribute.String("ens.query_kind", queryKind),
		attribute.String("result", result),
	))
}

func (m *EnsResolutionMetrics) RecordFallback(ctx context.Context, queryKind, outcome string, chainId uint64) {
	m.fallbackCounter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("ens.query_kind", queryKind),
		attribute.String("outcome", outcome),
		attribute.String("chain.id", strconv.FormatUint(chainId, 10)),
	))
}

func (m *EnsResolutionMetrics) RecordDecodeError(ctx context.Context, reason string) {
	m.decodeErrorCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("error.type", reason)))
}

func (m *EnsResolutionMetrics) RecordSignedResponse(ctx context.Context, err error) {
	outcome := "succeeded"
	if err != nil {
		outcome = "failed"
	}
	m.signedCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("outcome", outcome)))
}

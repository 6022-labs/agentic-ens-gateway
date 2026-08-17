package metrics

import (
	"context"
	"time"

	common_metrics "github.com/6022-labs/agentic-ens-gateway/src/common/metrics"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
)

const meterName = "agentic_ens_gateway_redis"

// RedisMetrics reports redis under the database semconv — redis is a db.system.name, not a bespoke cache tier.
type RedisMetrics struct {
	operationCounter  metric.Int64Counter
	durationHistogram metric.Float64Histogram
	errorCounter      metric.Int64Counter
}

func NewRedisMetrics() *RedisMetrics {
	noopMeter := noop.NewMeterProvider().Meter(meterName)
	meter := otel.GetMeterProvider().Meter(meterName)

	operationCounter, err := meter.Int64Counter(
		"db.client.operation.count",
		metric.WithDescription("Total number of redis commands"),
	)
	if err != nil {
		operationCounter, _ = noopMeter.Int64Counter("db.client.operation.count")
	}

	durationHistogram, err := meter.Float64Histogram(
		"db.client.operation.duration",
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(common_metrics.LatencySecondsBuckets...),
		metric.WithDescription("Duration of redis commands"),
	)
	if err != nil {
		durationHistogram, _ = noopMeter.Float64Histogram("db.client.operation.duration")
	}

	errorCounter, err := meter.Int64Counter(
		"db.client.operation.errors",
		metric.WithDescription("Total number of failed redis commands"),
	)
	if err != nil {
		errorCounter, _ = noopMeter.Int64Counter("db.client.operation.errors")
	}

	return &RedisMetrics{
		operationCounter:  operationCounter,
		durationHistogram: durationHistogram,
		errorCounter:      errorCounter,
	}
}

func (m *RedisMetrics) RecordOperation(ctx context.Context, operation string, duration time.Duration, failed bool) {
	attrs := metric.WithAttributes(
		semconv.DBSystemNameRedis,
		attribute.String("db.operation.name", operation),
	)
	m.operationCounter.Add(ctx, 1, attrs)
	m.durationHistogram.Record(ctx, duration.Seconds(), attrs)
	if failed {
		m.errorCounter.Add(ctx, 1, attrs)
	}
}

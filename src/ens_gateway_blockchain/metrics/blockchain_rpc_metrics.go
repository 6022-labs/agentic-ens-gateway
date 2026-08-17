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

const meterName = "agentic_ens_gateway_blockchain"

type BlockchainRpcMetrics struct {
	requestCounter    metric.Int64Counter
	durationHistogram metric.Float64Histogram
	errorCounter      metric.Int64Counter
}

func NewBlockchainRpcMetrics() *BlockchainRpcMetrics {
	noopMeter := noop.NewMeterProvider().Meter(meterName)
	meter := otel.GetMeterProvider().Meter(meterName)

	requestCounter, err := meter.Int64Counter(
		"rpc.request.count",
		metric.WithDescription("Total number of blockchain JSON-RPC requests"),
	)
	if err != nil {
		requestCounter, _ = noopMeter.Int64Counter("rpc.request.count")
	}

	durationHistogram, err := meter.Float64Histogram(
		"rpc.client.duration",
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(common_metrics.LatencySecondsBuckets...),
		metric.WithDescription("Duration of blockchain JSON-RPC requests"),
	)
	if err != nil {
		durationHistogram, _ = noopMeter.Float64Histogram("rpc.client.duration")
	}

	errorCounter, err := meter.Int64Counter(
		"rpc.request.errors",
		metric.WithDescription("Total number of failed blockchain JSON-RPC requests"),
	)
	if err != nil {
		errorCounter, _ = noopMeter.Int64Counter("rpc.request.errors")
	}

	return &BlockchainRpcMetrics{
		requestCounter:    requestCounter,
		durationHistogram: durationHistogram,
		errorCounter:      errorCounter,
	}
}

// RecordRequest records one JSON-RPC call: count, duration, and (when failed) an error.
func (m *BlockchainRpcMetrics) RecordRequest(ctx context.Context, chainID string, duration time.Duration, failed bool) {
	attrs := metric.WithAttributes(
		semconv.RPCSystemKey.String("jsonrpc"),
		attribute.String("chain.id", chainID),
	)
	m.requestCounter.Add(ctx, 1, attrs)
	m.durationHistogram.Record(ctx, duration.Seconds(), attrs)
	if failed {
		m.errorCounter.Add(ctx, 1, attrs)
	}
}

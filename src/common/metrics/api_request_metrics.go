package metrics

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

type ApiRequestMetrics struct {
	requestCounter    metric.Int64Counter
	durationHistogram metric.Float64Histogram
	activeRequests    metric.Int64UpDownCounter
	errorCounter      metric.Int64Counter
}

// NewApiRequestMetrics builds the http.server.* instruments under the given meter scope.
func NewApiRequestMetrics(meterName string) *ApiRequestMetrics {
	noopMeter := noop.NewMeterProvider().Meter(meterName)
	meter := otel.GetMeterProvider().Meter(meterName)

	requestCounter, err := meter.Int64Counter(
		"http.server.request.count",
		metric.WithDescription("Total number of HTTP requests handled by the server"),
	)
	if err != nil {
		requestCounter, _ = noopMeter.Int64Counter("http.server.request.count")
	}

	durationHistogram, err := meter.Float64Histogram(
		"http.server.request.duration",
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(LatencySecondsBuckets...),
		metric.WithDescription("Duration of HTTP requests"),
	)
	if err != nil {
		durationHistogram, _ = noopMeter.Float64Histogram("http.server.request.duration")
	}

	activeRequests, err := meter.Int64UpDownCounter(
		"http.server.active_requests",
		metric.WithDescription("Number of in-flight HTTP requests"),
	)
	if err != nil {
		activeRequests, _ = noopMeter.Int64UpDownCounter("http.server.active_requests")
	}

	errorCounter, err := meter.Int64Counter(
		"http.server.request.errors",
		metric.WithDescription("Total number of HTTP requests resulting in 5xx"),
	)
	if err != nil {
		errorCounter, _ = noopMeter.Int64Counter("http.server.request.errors")
	}

	return &ApiRequestMetrics{
		requestCounter:    requestCounter,
		durationHistogram: durationHistogram,
		activeRequests:    activeRequests,
		errorCounter:      errorCounter,
	}
}

func (m *ApiRequestMetrics) IncActiveRequests(ctx context.Context, method string) {
	m.activeRequests.Add(ctx, 1, metric.WithAttributes(attribute.String("http.request.method", method)))
}

func (m *ApiRequestMetrics) DecActiveRequests(ctx context.Context, method string) {
	m.activeRequests.Add(ctx, -1, metric.WithAttributes(attribute.String("http.request.method", method)))
}

// RecordRequest records one completed request: count, duration, and a 5xx error.
func (m *ApiRequestMetrics) RecordRequest(
	ctx context.Context,
	method, route string,
	statusCode int,
	duration time.Duration,
) {
	attrs := metric.WithAttributes(
		attribute.String("http.request.method", method),
		attribute.String("http.route", route),
		attribute.Int("http.response.status_code", statusCode),
	)
	m.requestCounter.Add(ctx, 1, attrs)
	m.durationHistogram.Record(ctx, duration.Seconds(), attrs)
	if statusCode >= 500 {
		m.errorCounter.Add(ctx, 1, attrs)
	}
}

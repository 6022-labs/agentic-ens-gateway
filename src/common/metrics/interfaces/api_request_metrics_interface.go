package interfaces

import (
	"context"
	"time"
)

// ApiRequestMetricsInterface records http.server.* metrics for an inbound HTTP server.
type ApiRequestMetricsInterface interface {
	IncActiveRequests(ctx context.Context, method string)
	DecActiveRequests(ctx context.Context, method string)
	RecordRequest(ctx context.Context, method, route string, statusCode int, duration time.Duration)
}

package middlewares

import (
	"time"

	metrics_interfaces "github.com/6022-labs/agentic-ens-gateway/src/common/metrics/interfaces"
	"github.com/gofiber/fiber/v2"
)

// ApiRequestMetricsMiddleware records http.server.* metrics around every request; must run before handlers.
type ApiRequestMetricsMiddleware struct {
	metrics metrics_interfaces.ApiRequestMetricsInterface
}

func NewApiRequestMetricsMiddleware(
	metrics metrics_interfaces.ApiRequestMetricsInterface,
) *ApiRequestMetricsMiddleware {
	return &ApiRequestMetricsMiddleware{
		metrics: metrics,
	}
}

func (m *ApiRequestMetricsMiddleware) Handle(ctx *fiber.Ctx) error {
	start := time.Now()
	method := ctx.Method()

	userCtx := ctx.UserContext()
	m.metrics.IncActiveRequests(userCtx, method)
	err := ctx.Next()
	m.metrics.DecActiveRequests(userCtx, method)
	if err != nil {
		_ = ctx.App().ErrorHandler(ctx, err)
		err = nil
	}

	// Prefer the route template (low cardinality) over the raw path.
	route := ctx.Path()
	if ctx.Route() != nil && ctx.Route().Path != "" {
		route = ctx.Route().Path
	}

	m.metrics.RecordRequest(userCtx, method, route, ctx.Response().StatusCode(), time.Since(start))

	return err
}

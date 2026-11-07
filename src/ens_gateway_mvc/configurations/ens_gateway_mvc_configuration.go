package configurations

import (
	common_metrics "github.com/6022-labs/agentic-ens-gateway/src/common/metrics"
	metrics_interfaces "github.com/6022-labs/agentic-ens-gateway/src/common/metrics/interfaces"
	"github.com/6022-labs/agentic-ens-gateway/src/common/mvc"
	mvc_interfaces "github.com/6022-labs/agentic-ens-gateway/src/common/mvc/interfaces"
	mvc_middlewares "github.com/6022-labs/agentic-ens-gateway/src/common/mvc/middlewares"
	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_mvc"
	"go.uber.org/dig"
)

func AddEnsGatewayMvcConfiguration(container *dig.Container) {
	if err := container.Provide(func() metrics_interfaces.ApiRequestMetricsInterface {
		return common_metrics.NewApiRequestMetrics("agentic_ens_gateway_api")
	}); err != nil {
		panic(err)
	}

	if err := container.Provide(mvc_middlewares.NewApiRequestMetricsMiddleware); err != nil {
		panic(err)
	}

	if err := container.Provide(mvc.NewErrorWriter); err != nil {
		panic(err)
	}

	// Static /health must register before /:sender/:data so Fiber matches it first.
	provideController(container, ens_gateway_mvc.NewHealthController)
	provideController(container, ens_gateway_mvc.NewCcipReadController)
}

func provideController(container *dig.Container, ctor any) {
	if err := container.Provide(
		ctor,
		dig.Group("controllers"),
		dig.As(new(mvc_interfaces.ControllerInterface)),
	); err != nil {
		panic(err)
	}
}

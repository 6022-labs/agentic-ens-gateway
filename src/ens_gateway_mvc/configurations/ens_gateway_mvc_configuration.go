package configurations

import (
	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_mvc"
	"go.uber.org/dig"
)

func AddEnsGatewayMvcConfiguration(container *dig.Container) {
	// HealthController is registered first so its static /health route is
	// matched ahead of CcipReadController's /:sender/:data param route.
	provideController(container, ens_gateway_mvc.NewHealthController)
	provideController(container, ens_gateway_mvc.NewCcipReadController)
}

func provideController(container *dig.Container, ctor any) {
	if err := container.Provide(
		ctor,
		dig.Group("controllers"),
		dig.As(new(ens_gateway_mvc.ControllerInterface)),
	); err != nil {
		panic(err)
	}
}

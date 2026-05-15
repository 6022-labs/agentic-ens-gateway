package configurations

import (
	ens_gateway_configurations "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/configurations"
	blockchain_configurations "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_blockchain/configurations"
	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_host/settings"
	mvc_configurations "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_mvc/configurations"
	"github.com/knadh/koanf/v2"
	"go.uber.org/dig"
)

func ConfigureDI(config *koanf.Koanf) *dig.Container {
	container := dig.New()

	if err := container.Provide(settings.NewHostSettings); err != nil {
		panic(err)
	}
	if err := container.Provide(func() *koanf.Koanf {
		return config
	}); err != nil {
		panic(err)
	}

	blockchain_configurations.AddEnsGatewayBlockchainConfiguration(container)
	ens_gateway_configurations.AddEnsGatewayConfiguration(container)
	mvc_configurations.AddEnsGatewayMvcConfiguration(container)

	return container
}

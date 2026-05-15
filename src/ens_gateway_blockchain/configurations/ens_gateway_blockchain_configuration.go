package configurations

import (
	ens_gateway_interfaces "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/services/interfaces"
	blockchain_services "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_blockchain/services"
	blockchain_settings "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_blockchain/settings"
	"go.uber.org/dig"
)

func AddEnsGatewayBlockchainConfiguration(container *dig.Container) {
	if err := container.Provide(blockchain_settings.NewChainsSettings); err != nil {
		panic(err)
	}

	if err := container.Provide(blockchain_settings.NewAgentEnsRegistriesSettings); err != nil {
		panic(err)
	}

	if err := container.Provide(
		blockchain_services.NewAgentEnsRegistryReader,
		dig.As(new(ens_gateway_interfaces.AgentEnsRegistryReaderInterface)),
	); err != nil {
		panic(err)
	}
}

package configurations

import (
	"github.com/6022-labs/agentic-ens-gateway/src/common/settings"
	ens_gateway_configurations "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/configurations"
	blockchain_configurations "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_blockchain/configurations"
	mvc_configurations "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_mvc/configurations"
	redis_configurations "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_redis/configurations"
	"github.com/knadh/koanf/v2"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

const (
	hostSettingsKey = "host"
	defaultApiPort  = 5000
)

func ConfigureDI(config *koanf.Koanf) *dig.Container {
	container := dig.New()

	if err := container.Provide(func(logger *zap.Logger, k *koanf.Koanf) *settings.HostSettings {
		return settings.NewHostSettings(logger, k, hostSettingsKey, defaultApiPort)
	}); err != nil {
		panic(err)
	}
	if err := container.Provide(func() *koanf.Koanf {
		return config
	}); err != nil {
		panic(err)
	}

	blockchain_configurations.AddEnsGatewayBlockchainConfiguration(container)
	redis_configurations.AddEnsGatewayRedisConfiguration(container)
	ens_gateway_configurations.AddEnsGatewayConfiguration(container)
	mvc_configurations.AddEnsGatewayMvcConfiguration(container)

	return container
}

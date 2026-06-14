package configurations

import (
	"strings"

	ens_gateway_services "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/services"
	ens_gateway_interfaces "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/services/interfaces"
	ens_gateway_settings "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/settings"
	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/use_cases"
	"github.com/ethereum/go-ethereum/crypto"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func AddEnsGatewayConfiguration(container *dig.Container) {
	if err := container.Provide(ens_gateway_settings.NewSignedResponseSettings); err != nil {
		panic(err)
	}

	if err := container.Provide(ens_gateway_settings.NewCacheSettings); err != nil {
		panic(err)
	}

	if err := container.Provide(
		newSigner,
		dig.As(new(ens_gateway_interfaces.SignerInterface)),
	); err != nil {
		panic(err)
	}

	if err := container.Provide(
		ens_gateway_services.NewResolverServiceCodec,
		dig.As(new(ens_gateway_interfaces.ResolverServiceCodecInterface)),
	); err != nil {
		panic(err)
	}

	if err := container.Provide(
		ens_gateway_services.NewSignedResponseBuilder,
		dig.As(new(ens_gateway_interfaces.ResponseBuilderInterface)),
	); err != nil {
		panic(err)
	}

	if err := container.Provide(
		ens_gateway_services.NewResolutionCache,
		dig.As(new(ens_gateway_interfaces.ResolutionCacheInterface)),
	); err != nil {
		panic(err)
	}

	if err := container.Provide(
		use_cases.NewResolveOffchainQuery,
		dig.As(new(use_cases.ResolveOffchainQueryInterface)),
	); err != nil {
		panic(err)
	}
}

func newSigner(
	logger *zap.Logger,
	signedResponseSettings *ens_gateway_settings.SignedResponseSettings,
) *ens_gateway_services.PrivateKeySigner {
	key, err := crypto.HexToECDSA(strings.TrimPrefix(signedResponseSettings.PrivateKey, "0x"))
	if err != nil {
		logger.Fatal("failed to parse signed_response.private_key", zap.Error(err))
	}

	signer := ens_gateway_services.NewPrivateKeySigner(key)
	logger.Info("gateway signer loaded", zap.String("address", signer.Address().Hex()))
	return signer
}

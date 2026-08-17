package services

import (
	"context"
	"fmt"

	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/services/interfaces"
	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_blockchain/abi"
	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_blockchain/factory"
	blockchain_settings "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_blockchain/settings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

type UniversalResolverClient struct {
	resolverServiceCodec interfaces.ResolverServiceCodecInterface
	bindings             map[uint64]*abi.UniversalResolver
}

func NewUniversalResolverClient(
	logger *zap.Logger,
	chainsSettings *blockchain_settings.ChainsSettings,
	universalResolverSettings *blockchain_settings.UniversalResolverSettings,
	resolverServiceCodec interfaces.ResolverServiceCodecInterface,
	ethClientFactory *factory.EthClientFactory,
) *UniversalResolverClient {
	bindings := map[uint64]*abi.UniversalResolver{}

	for sourceChainId := range universalResolverSettings.Mappings {
		resolverChainId, addr, ok := universalResolverSettings.Get(sourceChainId)
		if !ok {
			logger.Info("ENS fallback disabled for source chain: no resolver address for mapped chain",
				zap.Uint64("source_chain_id", sourceChainId),
				zap.Uint64("universal_resolver_chain_id", resolverChainId),
			)
			continue
		}

		resolverChain := chainsSettings.Get(resolverChainId)
		if resolverChain == nil {
			logger.Info("ENS fallback disabled for source chain: add chains.<universal_resolver_chain_id>.rpc_http_url to enable",
				zap.Uint64("source_chain_id", sourceChainId),
				zap.Uint64("universal_resolver_chain_id", resolverChainId),
			)
			continue
		}

		client, err := ethClientFactory.For(resolverChainId)
		if err != nil {
			logger.Fatal("failed to dial RPC for ENS fallback",
				zap.Uint64("universal_resolver_chain_id", resolverChainId),
				zap.Error(err),
			)
		}

		binding, err := abi.NewUniversalResolver(addr, client)
		if err != nil {
			logger.Fatal("failed to bind UniversalResolver",
				zap.Uint64("universal_resolver_chain_id", resolverChainId),
				zap.String("address", addr.Hex()),
				zap.Error(err),
			)
		}
		bindings[sourceChainId] = binding

		logger.Info("ENS fallback enabled",
			zap.Uint64("source_chain_id", sourceChainId),
			zap.Uint64("universal_resolver_chain_id", resolverChainId),
			zap.String("universal_resolver", addr.Hex()),
		)
	}

	return &UniversalResolverClient{resolverServiceCodec: resolverServiceCodec, bindings: bindings}
}

func (c *UniversalResolverClient) ResolveText(ctx context.Context, sourceChainId uint64, dnsName, innerData []byte) (string, error) {
	result, ok, err := c.resolve(ctx, sourceChainId, dnsName, innerData)
	if err != nil || !ok {
		return "", err
	}
	return c.resolverServiceCodec.DecodeTextResult(result)
}

func (c *UniversalResolverClient) ResolveAddr(ctx context.Context, sourceChainId uint64, dnsName, innerData []byte) (common.Address, error) {
	result, ok, err := c.resolve(ctx, sourceChainId, dnsName, innerData)
	if err != nil || !ok {
		return common.Address{}, err
	}
	return c.resolverServiceCodec.DecodeAddrResult(result)
}

func (c *UniversalResolverClient) resolve(ctx context.Context, sourceChainId uint64, dnsName, innerData []byte) ([]byte, bool, error) {
	binding, ok := c.bindings[sourceChainId]
	if !ok {
		return nil, false, nil
	}
	result, _, err := binding.Resolve(&bind.CallOpts{Context: ctx}, dnsName, innerData)
	if err != nil {
		return nil, false, fmt.Errorf("universal resolver resolve: %w", err)
	}
	return result, true, nil
}

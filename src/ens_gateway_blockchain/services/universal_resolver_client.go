package services

import (
	"context"
	"fmt"

	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/services/interfaces"
	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_blockchain/abi"
	blockchain_settings "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_blockchain/settings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

type UniversalResolverClient struct {
	codec    interfaces.ResolverServiceCodecInterface
	bindings map[uint64]*abi.UniversalResolver
}

func NewUniversalResolverClient(
	logger *zap.Logger,
	chains *blockchain_settings.ChainsSettings,
	fallbacks *blockchain_settings.UniversalResolverSettings,
	codec interfaces.ResolverServiceCodecInterface,
) *UniversalResolverClient {
	resolverClients := map[uint64]*ethclient.Client{}
	bindings := map[uint64]*abi.UniversalResolver{}

	for sourceChainId := range fallbacks.Mappings {
		resolverChainId, addr, ok := fallbacks.Get(sourceChainId)
		if !ok {
			logger.Info("ENS fallback disabled for source chain: no resolver address for mapped chain",
				zap.Uint64("source_chain_id", sourceChainId),
				zap.Uint64("universal_resolver_chain_id", resolverChainId),
			)
			continue
		}

		resolverChain := chains.Get(resolverChainId)
		if resolverChain == nil {
			logger.Info("ENS fallback disabled for source chain: add chains.<universal_resolver_chain_id>.rpc_http_url to enable",
				zap.Uint64("source_chain_id", sourceChainId),
				zap.Uint64("universal_resolver_chain_id", resolverChainId),
			)
			continue
		}

		client, ok := resolverClients[resolverChainId]
		if !ok {
			c, err := ethclient.Dial(resolverChain.RpcHttpUrl)
			if err != nil {
				logger.Fatal("failed to dial RPC for ENS fallback",
					zap.Uint64("universal_resolver_chain_id", resolverChainId),
					zap.Error(err),
				)
			}
			resolverClients[resolverChainId] = c
			client = c
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

	return &UniversalResolverClient{codec: codec, bindings: bindings}
}

func (c *UniversalResolverClient) Resolve(ctx context.Context, sourceChainId uint64, dnsName, innerData []byte) (string, error) {
	binding, ok := c.bindings[sourceChainId]
	if !ok {
		return "", nil
	}

	result, _, err := binding.Resolve(&bind.CallOpts{Context: ctx}, dnsName, innerData)
	if err != nil {
		return "", fmt.Errorf("universal resolver resolve: %w", err)
	}

	return c.codec.DecodeTextResult(result)
}

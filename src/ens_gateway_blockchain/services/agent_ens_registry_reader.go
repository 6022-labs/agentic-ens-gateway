package services

import (
	"context"

	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_blockchain/abi"
	blockchain_settings "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_blockchain/settings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

type AgentEnsRegistryReader struct {
	bindings map[uint64]*abi.AgentEnsRegistry
}

// Iterates registries (not chains); extra chains in ChainsSettings are ignored.
func NewAgentEnsRegistryReader(
	logger *zap.Logger,
	chains *blockchain_settings.ChainsSettings,
	registries *blockchain_settings.AgentEnsRegistriesSettings,
) *AgentEnsRegistryReader {
	bindings := make(map[uint64]*abi.AgentEnsRegistry, len(*registries))

	for chainId, address := range *registries {
		chain := chains.Get(chainId)
		if chain == nil {
			logger.Fatal(
				"AgentEnsRegistry configured for a chain with no RPC url; add the chain to appsettings.json (chains.<id>.rpc_http_url)",
				zap.Uint64("chain_id", chainId),
			)
		}

		client, err := ethclient.Dial(chain.RpcHttpUrl)
		if err != nil {
			logger.Fatal("failed to dial EVM RPC",
				zap.Uint64("chain_id", chainId),
				zap.Error(err),
			)
		}

		binding, err := abi.NewAgentEnsRegistry(address, client)
		if err != nil {
			logger.Fatal("failed to bind AgentEnsRegistry",
				zap.Uint64("chain_id", chainId),
				zap.Error(err),
			)
		}

		bindings[chainId] = binding
		logger.Info("bound AgentEnsRegistry",
			zap.Uint64("chain_id", chainId),
			zap.String("address", address.Hex()),
		)
	}

	return &AgentEnsRegistryReader{bindings: bindings}
}

// Unbound chain id returns ("", nil) so the caller can fall back to L1.
func (r *AgentEnsRegistryReader) Text(ctx context.Context, chainId uint64, node [32]byte, key string) (string, error) {
	binding, ok := r.bindings[chainId]
	if !ok {
		return "", nil
	}
	return binding.Text(&bind.CallOpts{Context: ctx}, node, key)
}

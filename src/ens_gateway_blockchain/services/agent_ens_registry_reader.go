package services

import (
	"context"
	"fmt"

	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_blockchain/abi"
	blockchain_settings "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_blockchain/settings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

type AgentEnsRegistryReader struct {
	bindings map[uint64]*abi.AgentEnsRegistry
}

func NewAgentEnsRegistryReader(
	logger *zap.Logger,
	chains *blockchain_settings.ChainsSettings,
	registries *blockchain_settings.AgentEnsRegistriesSettings,
) *AgentEnsRegistryReader {
	bindings := make(map[uint64]*abi.AgentEnsRegistry, len(*chains))

	for chainId, chain := range *chains {
		address, ok := registries.Get(chainId)
		if !ok {
			logger.Fatal(
				"no AgentEnsRegistry address for configured chain; add an entry to agent_ens_registries.json and rebuild",
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

func (r *AgentEnsRegistryReader) Text(ctx context.Context, chainId uint64, node [32]byte, key string) (string, error) {
	binding, ok := r.bindings[chainId]
	if !ok {
		return "", fmt.Errorf("unsupported chain id: %d", chainId)
	}
	return binding.Text(&bind.CallOpts{Context: ctx}, node, key)
}

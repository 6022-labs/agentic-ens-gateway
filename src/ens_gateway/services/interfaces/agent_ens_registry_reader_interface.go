package interfaces

import "context"

// AgentEnsRegistryReaderInterface is the outbound port for reading agent text
// records from the AgentEnsRegistry contract on the right Polygon-side chain.
// The implementation routes by `chainId` and returns an error for unconfigured
// chains.
type AgentEnsRegistryReaderInterface interface {
	Text(ctx context.Context, chainId uint64, node [32]byte, key string) (string, error)
}

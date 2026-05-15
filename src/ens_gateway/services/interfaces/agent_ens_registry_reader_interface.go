package interfaces

import "context"

type AgentEnsRegistryReaderInterface interface {
	Text(ctx context.Context, chainId uint64, node [32]byte, key string) (string, error)
}

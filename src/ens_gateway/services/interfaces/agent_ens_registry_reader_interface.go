package interfaces

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

type AgentEnsRegistryReaderInterface interface {
	Text(ctx context.Context, chainId uint64, node [32]byte, key string) (string, error)
	Addr(ctx context.Context, chainId uint64, node [32]byte) (common.Address, error)
}

package interfaces

import (
	"context"
	"time"
)

// BlockchainRpcMetricsInterface records JSON-RPC call metrics, shared by the on-chain adapters' eth client factories.
type BlockchainRpcMetricsInterface interface {
	RecordRequest(ctx context.Context, chainID string, duration time.Duration, failed bool)
}

package interfaces

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

// ResponseBuilderInterface builds the ERC-3668 lookup response a CCIP-Read
// client hands back to the L1 resolver's callback. The signed implementation is
// the current default; a proof-based implementation can be swapped in behind
// this port without touching the use case.
type ResponseBuilderInterface interface {
	Build(ctx context.Context, resolver common.Address, callData []byte, result []byte) ([]byte, error)
}

package interfaces

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

type ResponseBuilderInterface interface {
	Build(ctx context.Context, resolver common.Address, callData []byte, result []byte) ([]byte, error)
}

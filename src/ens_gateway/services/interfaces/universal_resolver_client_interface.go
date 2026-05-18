package interfaces

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

type UniversalResolverClientInterface interface {
	ResolveText(ctx context.Context, sourceChainId uint64, dnsName, innerData []byte) (string, error)
	ResolveAddr(ctx context.Context, sourceChainId uint64, dnsName, innerData []byte) (common.Address, error)
}

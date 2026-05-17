package interfaces

import "context"

type UniversalResolverClientInterface interface {
	Resolve(ctx context.Context, sourceChainId uint64, dnsName, innerData []byte) (string, error)
}

package interfaces

import "context"

// RecordResolverInterface answers one decoded ENS query, in cache → L2 registry → L1 fallback order.
type RecordResolverInterface interface {
	// Resolve returns the ABI-encoded result plus the tier that answered it.
	Resolve(ctx context.Context, query DecodedQuery) ([]byte, string, error)
}

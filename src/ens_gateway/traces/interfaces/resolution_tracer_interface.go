package interfaces

import "context"

// ResolutionTracerInterface opens the spans of one CCIP-Read resolution.
type ResolutionTracerInterface interface {
	// StartResolution opens the root span; queryKind is text or addr.
	StartResolution(ctx context.Context, queryKind string, chainId uint64) (context.Context, ResolutionSpan)
	StartRegistryRead(ctx context.Context, chainId uint64) (context.Context, Span)
	StartUniversalResolverFallback(ctx context.Context, chainId uint64) (context.Context, Span)
}

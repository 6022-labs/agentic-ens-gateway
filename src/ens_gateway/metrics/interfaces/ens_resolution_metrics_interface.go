package interfaces

import (
	"context"
	"time"
)

const (
	ResolutionSourceCache             = "cache"
	ResolutionSourceAgentEnsRegistry  = "agent_ens_registry"
	ResolutionSourceUniversalResolver = "universal_resolver"
)

const (
	FallbackOutcomeResolved = "resolved"
	FallbackOutcomeEmpty    = "empty"
	FallbackOutcomeFailed   = "failed"
)

// EnsResolutionMetricsInterface records the ens.* metrics of one CCIP-Read resolution.
type EnsResolutionMetricsInterface interface {
	// RecordResolution records one resolved query; source tells which tier answered it.
	RecordResolution(ctx context.Context, queryKind, source string, chainId uint64, duration time.Duration, err error)
	RecordCacheResult(ctx context.Context, queryKind string, hit bool)
	// RecordFallback records an L1 universal-resolver attempt, made only when the L2 registry answered empty.
	RecordFallback(ctx context.Context, queryKind, outcome string, chainId uint64)
	RecordDecodeError(ctx context.Context, reason string)
	RecordSignedResponse(ctx context.Context, err error)
}

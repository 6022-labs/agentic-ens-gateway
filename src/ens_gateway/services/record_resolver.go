package services

import (
	"context"
	"fmt"

	metrics_interfaces "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/metrics/interfaces"
	ens_gateway_interfaces "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/services/interfaces"
	traces_interfaces "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/traces/interfaces"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

// RecordResolver reads one ENS record, preferring the cache, then the L2 registry, then the L1 universal resolver.
type RecordResolver struct {
	logger                  *zap.Logger
	resolverServiceCodec    ens_gateway_interfaces.ResolverServiceCodecInterface
	agentEnsRegistryReader  ens_gateway_interfaces.AgentEnsRegistryReaderInterface
	universalResolverClient ens_gateway_interfaces.UniversalResolverClientInterface
	resolutionCache         ens_gateway_interfaces.ResolutionCacheInterface
	ensResolutionMetrics    metrics_interfaces.EnsResolutionMetricsInterface
	resolutionTracer        traces_interfaces.ResolutionTracerInterface
}

func NewRecordResolver(
	logger *zap.Logger,
	resolverServiceCodec ens_gateway_interfaces.ResolverServiceCodecInterface,
	agentEnsRegistryReader ens_gateway_interfaces.AgentEnsRegistryReaderInterface,
	universalResolverClient ens_gateway_interfaces.UniversalResolverClientInterface,
	resolutionCache ens_gateway_interfaces.ResolutionCacheInterface,
	ensResolutionMetrics metrics_interfaces.EnsResolutionMetricsInterface,
	resolutionTracer traces_interfaces.ResolutionTracerInterface,
) *RecordResolver {
	return &RecordResolver{
		logger:                  logger,
		resolverServiceCodec:    resolverServiceCodec,
		agentEnsRegistryReader:  agentEnsRegistryReader,
		universalResolverClient: universalResolverClient,
		resolutionCache:         resolutionCache,
		ensResolutionMetrics:    ensResolutionMetrics,
		resolutionTracer:        resolutionTracer,
	}
}

func (s *RecordResolver) Resolve(
	ctx context.Context,
	query ens_gateway_interfaces.DecodedQuery,
) ([]byte, string, error) {
	if query.Kind == ens_gateway_interfaces.QueryKindText {
		return s.resolveText(ctx, query)
	}
	return s.resolveAddr(ctx, query)
}

func (s *RecordResolver) resolveText(
	ctx context.Context,
	query ens_gateway_interfaces.DecodedQuery,
) ([]byte, string, error) {
	value, source, err := s.textValue(ctx, query)
	if err != nil {
		return nil, source, err
	}

	encoded, err := s.resolverServiceCodec.EncodeTextResult(value)
	if err != nil {
		return nil, source, fmt.Errorf("encode text result: %w", err)
	}
	return encoded, source, nil
}

func (s *RecordResolver) textValue(
	ctx context.Context,
	query ens_gateway_interfaces.DecodedQuery,
) (string, string, error) {
	cached := s.resolutionCache.GetText(query.ChainId, query.Node, query.Key)
	s.ensResolutionMetrics.RecordCacheResult(ctx, "text", cached != nil)
	if cached != nil {
		return *cached, metrics_interfaces.ResolutionSourceCache, nil
	}

	source := metrics_interfaces.ResolutionSourceAgentEnsRegistry
	readCtx, readSpan := s.resolutionTracer.StartRegistryRead(ctx, query.ChainId)
	value, err := s.agentEnsRegistryReader.Text(readCtx, query.ChainId, query.Node, query.Key)
	if err != nil {
		readSpan.Fail(err)
		readSpan.End()
		return "", source, fmt.Errorf("read agent ens registry: %w", err)
	}
	readSpan.End()
	if value == "" {
		fallbackCtx, fallbackSpan := s.resolutionTracer.StartUniversalResolverFallback(ctx, query.ChainId)
		fallback, fallbackErr := s.universalResolverClient.ResolveText(
			fallbackCtx, query.ChainId, query.DnsName, query.InnerData,
		)
		switch {
		case fallbackErr != nil:
			fallbackSpan.Fail(fallbackErr)
			s.ensResolutionMetrics.RecordFallback(ctx, "text", metrics_interfaces.FallbackOutcomeFailed, query.ChainId)
			s.logger.Warn("L1 universal resolver fallback failed",
				zap.Error(fallbackErr),
				zap.Uint64("chain_id", query.ChainId),
				zap.String("key", query.Key),
			)
		case fallback != "":
			s.ensResolutionMetrics.RecordFallback(
				ctx,
				"text",
				metrics_interfaces.FallbackOutcomeResolved,
				query.ChainId,
			)
			source = metrics_interfaces.ResolutionSourceUniversalResolver
			value = fallback
		default:
			s.ensResolutionMetrics.RecordFallback(ctx, "text", metrics_interfaces.FallbackOutcomeEmpty, query.ChainId)
		}
		fallbackSpan.End()
	}

	s.resolutionCache.SetText(query.ChainId, query.Node, query.Key, value)
	return value, source, nil
}

func (s *RecordResolver) resolveAddr(
	ctx context.Context,
	query ens_gateway_interfaces.DecodedQuery,
) ([]byte, string, error) {
	value, source, err := s.addrValue(ctx, query)
	if err != nil {
		return nil, source, err
	}

	encoded, err := s.resolverServiceCodec.EncodeAddrResult(value)
	if err != nil {
		return nil, source, fmt.Errorf("encode addr result: %w", err)
	}
	return encoded, source, nil
}

func (s *RecordResolver) addrValue(
	ctx context.Context,
	query ens_gateway_interfaces.DecodedQuery,
) (common.Address, string, error) {
	cached := s.resolutionCache.GetAddr(query.ChainId, query.Node)
	s.ensResolutionMetrics.RecordCacheResult(ctx, "addr", cached != nil)
	if cached != nil {
		return *cached, metrics_interfaces.ResolutionSourceCache, nil
	}

	source := metrics_interfaces.ResolutionSourceAgentEnsRegistry
	readCtx, readSpan := s.resolutionTracer.StartRegistryRead(ctx, query.ChainId)
	value, err := s.agentEnsRegistryReader.Addr(readCtx, query.ChainId, query.Node)
	if err != nil {
		readSpan.Fail(err)
		readSpan.End()
		return common.Address{}, source, fmt.Errorf("read agent ens registry: %w", err)
	}
	readSpan.End()
	if value == (common.Address{}) {
		fallbackCtx, fallbackSpan := s.resolutionTracer.StartUniversalResolverFallback(ctx, query.ChainId)
		fallback, fallbackErr := s.universalResolverClient.ResolveAddr(
			fallbackCtx, query.ChainId, query.DnsName, query.InnerData,
		)
		switch {
		case fallbackErr != nil:
			fallbackSpan.Fail(fallbackErr)
			s.ensResolutionMetrics.RecordFallback(ctx, "addr", metrics_interfaces.FallbackOutcomeFailed, query.ChainId)
			s.logger.Warn("L1 universal resolver fallback failed",
				zap.Error(fallbackErr),
				zap.Uint64("chain_id", query.ChainId),
			)
		case fallback != (common.Address{}):
			s.ensResolutionMetrics.RecordFallback(
				ctx,
				"addr",
				metrics_interfaces.FallbackOutcomeResolved,
				query.ChainId,
			)
			source = metrics_interfaces.ResolutionSourceUniversalResolver
			value = fallback
		default:
			s.ensResolutionMetrics.RecordFallback(ctx, "addr", metrics_interfaces.FallbackOutcomeEmpty, query.ChainId)
		}
		fallbackSpan.End()
	}

	s.resolutionCache.SetAddr(query.ChainId, query.Node, value)
	return value, source, nil
}

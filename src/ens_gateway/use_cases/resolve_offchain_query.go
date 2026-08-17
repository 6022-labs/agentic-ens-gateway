package use_cases

import (
	"context"
	"fmt"
	"time"

	"github.com/6022-labs/agentic-ens-gateway/src/common/errors"
	metrics_interfaces "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/metrics/interfaces"
	ens_gateway_interfaces "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/services/interfaces"
	traces_interfaces "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/traces/interfaces"
	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/use_cases/responses"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"go.uber.org/zap"
)

type ResolveOffchainQueryInterface interface {
	Execute(ctx context.Context, request *requests.ResolveOffchainQueryRequest) (*responses.CcipReadResponse, error)
}

type ResolveOffchainQuery struct {
	logger                  *zap.Logger
	resolverServiceCodec    ens_gateway_interfaces.ResolverServiceCodecInterface
	agentEnsRegistryReader  ens_gateway_interfaces.AgentEnsRegistryReaderInterface
	universalResolverClient ens_gateway_interfaces.UniversalResolverClientInterface
	responseBuilder         ens_gateway_interfaces.ResponseBuilderInterface
	resolutionCache         ens_gateway_interfaces.ResolutionCacheInterface
	ensResolutionMetrics    metrics_interfaces.EnsResolutionMetricsInterface
	resolutionTracer        traces_interfaces.ResolutionTracerInterface
}

func NewResolveOffchainQuery(
	logger *zap.Logger,
	resolverServiceCodec ens_gateway_interfaces.ResolverServiceCodecInterface,
	agentEnsRegistryReader ens_gateway_interfaces.AgentEnsRegistryReaderInterface,
	universalResolverClient ens_gateway_interfaces.UniversalResolverClientInterface,
	responseBuilder ens_gateway_interfaces.ResponseBuilderInterface,
	resolutionCache ens_gateway_interfaces.ResolutionCacheInterface,
	ensResolutionMetrics metrics_interfaces.EnsResolutionMetricsInterface,
	resolutionTracer traces_interfaces.ResolutionTracerInterface,
) *ResolveOffchainQuery {
	return &ResolveOffchainQuery{
		logger:                  logger,
		resolverServiceCodec:    resolverServiceCodec,
		agentEnsRegistryReader:  agentEnsRegistryReader,
		universalResolverClient: universalResolverClient,
		responseBuilder:         responseBuilder,
		resolutionCache:         resolutionCache,
		ensResolutionMetrics:    ensResolutionMetrics,
		resolutionTracer:        resolutionTracer,
	}
}

func (u *ResolveOffchainQuery) Execute(
	ctx context.Context,
	request *requests.ResolveOffchainQueryRequest,
) (*responses.CcipReadResponse, error) {
	query, err := u.resolverServiceCodec.Decode(request.Data)
	if err != nil {
		u.ensResolutionMetrics.RecordDecodeError(ctx, "undecodable_calldata")
		return nil, errors.NewValidationError("data", fmt.Sprintf("invalid resolver call data: %s", err.Error()))
	}

	if query.Kind != ens_gateway_interfaces.QueryKindText && query.Kind != ens_gateway_interfaces.QueryKindAddr {
		u.ensResolutionMetrics.RecordDecodeError(ctx, "unsupported_query_kind")
		return nil, errors.NewValidationError("data", fmt.Sprintf("unsupported query kind: %d", query.Kind))
	}

	queryKind := u.queryKindName(query.Kind)
	start := time.Now()
	ctx, span := u.resolutionTracer.StartResolution(ctx, queryKind, query.ChainId)
	defer span.End()

	var result []byte
	var source string
	if query.Kind == ens_gateway_interfaces.QueryKindText {
		result, source, err = u.resolveText(ctx, query)
	} else {
		result, source, err = u.resolveAddr(ctx, query)
	}

	span.RecordSource(source)
	u.ensResolutionMetrics.RecordResolution(ctx, queryKind, source, query.ChainId, time.Since(start), err)
	if err != nil {
		span.Fail(err)
		return nil, err
	}

	signed, err := u.responseBuilder.Build(ctx, request.Sender, request.Data, result)
	u.ensResolutionMetrics.RecordSignedResponse(ctx, err)
	if err != nil {
		span.Fail(err)
		return nil, fmt.Errorf("build signed response: %w", err)
	}

	return &responses.CcipReadResponse{Data: hexutil.Encode(signed)}, nil
}

func (u *ResolveOffchainQuery) queryKindName(kind ens_gateway_interfaces.QueryKind) string {
	if kind == ens_gateway_interfaces.QueryKindAddr {
		return "addr"
	}
	return "text"
}

func (u *ResolveOffchainQuery) resolveText(
	ctx context.Context,
	query ens_gateway_interfaces.DecodedQuery,
) ([]byte, string, error) {
	value, source, err := u.textValue(ctx, query)
	if err != nil {
		return nil, source, err
	}

	encoded, err := u.resolverServiceCodec.EncodeTextResult(value)
	if err != nil {
		return nil, source, fmt.Errorf("encode text result: %w", err)
	}
	return encoded, source, nil
}

func (u *ResolveOffchainQuery) textValue(
	ctx context.Context,
	query ens_gateway_interfaces.DecodedQuery,
) (string, string, error) {
	cached := u.resolutionCache.GetText(query.ChainId, query.Node, query.Key)
	u.ensResolutionMetrics.RecordCacheResult(ctx, "text", cached != nil)
	if cached != nil {
		return *cached, metrics_interfaces.ResolutionSourceCache, nil
	}

	source := metrics_interfaces.ResolutionSourceAgentEnsRegistry
	readCtx, readSpan := u.resolutionTracer.StartRegistryRead(ctx, query.ChainId)
	value, err := u.agentEnsRegistryReader.Text(readCtx, query.ChainId, query.Node, query.Key)
	if err != nil {
		readSpan.Fail(err)
		readSpan.End()
		return "", source, fmt.Errorf("read agent ens registry: %w", err)
	}
	readSpan.End()
	if value == "" {
		fallbackCtx, fallbackSpan := u.resolutionTracer.StartUniversalResolverFallback(ctx, query.ChainId)
		fallback, fallbackErr := u.universalResolverClient.ResolveText(
			fallbackCtx, query.ChainId, query.DnsName, query.InnerData,
		)
		switch {
		case fallbackErr != nil:
			fallbackSpan.Fail(fallbackErr)
			u.ensResolutionMetrics.RecordFallback(ctx, "text", metrics_interfaces.FallbackOutcomeFailed, query.ChainId)
			u.logger.Warn("L1 universal resolver fallback failed",
				zap.Error(fallbackErr),
				zap.Uint64("chain_id", query.ChainId),
				zap.String("key", query.Key),
			)
		case fallback != "":
			u.ensResolutionMetrics.RecordFallback(ctx, "text", metrics_interfaces.FallbackOutcomeResolved, query.ChainId)
			source = metrics_interfaces.ResolutionSourceUniversalResolver
			value = fallback
		default:
			u.ensResolutionMetrics.RecordFallback(ctx, "text", metrics_interfaces.FallbackOutcomeEmpty, query.ChainId)
		}
		fallbackSpan.End()
	}

	u.resolutionCache.SetText(query.ChainId, query.Node, query.Key, value)
	return value, source, nil
}

func (u *ResolveOffchainQuery) resolveAddr(
	ctx context.Context,
	query ens_gateway_interfaces.DecodedQuery,
) ([]byte, string, error) {
	value, source, err := u.addrValue(ctx, query)
	if err != nil {
		return nil, source, err
	}

	encoded, err := u.resolverServiceCodec.EncodeAddrResult(value)
	if err != nil {
		return nil, source, fmt.Errorf("encode addr result: %w", err)
	}
	return encoded, source, nil
}

func (u *ResolveOffchainQuery) addrValue(
	ctx context.Context,
	query ens_gateway_interfaces.DecodedQuery,
) (common.Address, string, error) {
	cached := u.resolutionCache.GetAddr(query.ChainId, query.Node)
	u.ensResolutionMetrics.RecordCacheResult(ctx, "addr", cached != nil)
	if cached != nil {
		return *cached, metrics_interfaces.ResolutionSourceCache, nil
	}

	source := metrics_interfaces.ResolutionSourceAgentEnsRegistry
	readCtx, readSpan := u.resolutionTracer.StartRegistryRead(ctx, query.ChainId)
	value, err := u.agentEnsRegistryReader.Addr(readCtx, query.ChainId, query.Node)
	if err != nil {
		readSpan.Fail(err)
		readSpan.End()
		return common.Address{}, source, fmt.Errorf("read agent ens registry: %w", err)
	}
	readSpan.End()
	if value == (common.Address{}) {
		fallbackCtx, fallbackSpan := u.resolutionTracer.StartUniversalResolverFallback(ctx, query.ChainId)
		fallback, fallbackErr := u.universalResolverClient.ResolveAddr(
			fallbackCtx, query.ChainId, query.DnsName, query.InnerData,
		)
		switch {
		case fallbackErr != nil:
			fallbackSpan.Fail(fallbackErr)
			u.ensResolutionMetrics.RecordFallback(ctx, "addr", metrics_interfaces.FallbackOutcomeFailed, query.ChainId)
			u.logger.Warn("L1 universal resolver fallback failed",
				zap.Error(fallbackErr),
				zap.Uint64("chain_id", query.ChainId),
			)
		case fallback != (common.Address{}):
			u.ensResolutionMetrics.RecordFallback(ctx, "addr", metrics_interfaces.FallbackOutcomeResolved, query.ChainId)
			source = metrics_interfaces.ResolutionSourceUniversalResolver
			value = fallback
		default:
			u.ensResolutionMetrics.RecordFallback(ctx, "addr", metrics_interfaces.FallbackOutcomeEmpty, query.ChainId)
		}
		fallbackSpan.End()
	}

	u.resolutionCache.SetAddr(query.ChainId, query.Node, value)
	return value, source, nil
}

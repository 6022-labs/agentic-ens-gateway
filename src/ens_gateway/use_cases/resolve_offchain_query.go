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
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type ResolveOffchainQuery struct {
	resolverServiceCodec ens_gateway_interfaces.ResolverServiceCodecInterface
	recordResolver       ens_gateway_interfaces.RecordResolverInterface
	responseBuilder      ens_gateway_interfaces.ResponseBuilderInterface
	ensResolutionMetrics metrics_interfaces.EnsResolutionMetricsInterface
	resolutionTracer     traces_interfaces.ResolutionTracerInterface
}

func NewResolveOffchainQuery(
	resolverServiceCodec ens_gateway_interfaces.ResolverServiceCodecInterface,
	recordResolver ens_gateway_interfaces.RecordResolverInterface,
	responseBuilder ens_gateway_interfaces.ResponseBuilderInterface,
	ensResolutionMetrics metrics_interfaces.EnsResolutionMetricsInterface,
	resolutionTracer traces_interfaces.ResolutionTracerInterface,
) *ResolveOffchainQuery {
	return &ResolveOffchainQuery{
		resolverServiceCodec: resolverServiceCodec,
		recordResolver:       recordResolver,
		responseBuilder:      responseBuilder,
		ensResolutionMetrics: ensResolutionMetrics,
		resolutionTracer:     resolutionTracer,
	}
}

func (u *ResolveOffchainQuery) Execute(
	ctx context.Context,
	request *requests.ResolveOffchainQueryRequest,
) (*responses.CcipReadResponse, error) {
	if err := request.ValidateAndSanitize(); err != nil {
		return nil, err
	}

	callData := request.CallData()

	query, err := u.resolverServiceCodec.Decode(callData)
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

	result, source, err := u.recordResolver.Resolve(ctx, query)

	span.RecordSource(source)
	u.ensResolutionMetrics.RecordResolution(ctx, queryKind, source, query.ChainId, time.Since(start), err)
	if err != nil {
		span.Fail(err)
		return nil, errors.NewInternalError("internal_error", fmt.Errorf("resolve record: %w", err))
	}

	signed, err := u.responseBuilder.Build(ctx, request.SenderAddress(), callData, result)
	u.ensResolutionMetrics.RecordSignedResponse(ctx, err)
	if err != nil {
		span.Fail(err)
		return nil, errors.NewInternalError("internal_error", fmt.Errorf("build signed response: %w", err))
	}

	return &responses.CcipReadResponse{Data: hexutil.Encode(signed)}, nil
}

func (u *ResolveOffchainQuery) queryKindName(kind ens_gateway_interfaces.QueryKind) string {
	if kind == ens_gateway_interfaces.QueryKindAddr {
		return "addr"
	}
	return "text"
}

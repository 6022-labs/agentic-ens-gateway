package use_cases

import (
	"context"
	"fmt"

	ens_gateway_interfaces "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/services/interfaces"
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
}

func NewResolveOffchainQuery(
	logger *zap.Logger,
	resolverServiceCodec ens_gateway_interfaces.ResolverServiceCodecInterface,
	agentEnsRegistryReader ens_gateway_interfaces.AgentEnsRegistryReaderInterface,
	universalResolverClient ens_gateway_interfaces.UniversalResolverClientInterface,
	responseBuilder ens_gateway_interfaces.ResponseBuilderInterface,
	resolutionCache ens_gateway_interfaces.ResolutionCacheInterface,
) *ResolveOffchainQuery {
	return &ResolveOffchainQuery{
		logger:                  logger,
		resolverServiceCodec:    resolverServiceCodec,
		agentEnsRegistryReader:  agentEnsRegistryReader,
		universalResolverClient: universalResolverClient,
		responseBuilder:         responseBuilder,
		resolutionCache:         resolutionCache,
	}
}

func (u *ResolveOffchainQuery) Execute(
	ctx context.Context,
	request *requests.ResolveOffchainQueryRequest,
) (*responses.CcipReadResponse, error) {
	query, err := u.resolverServiceCodec.Decode(request.Data)
	if err != nil {
		return nil, fmt.Errorf("decode call data: %w", err)
	}

	var result []byte
	switch query.Kind {
	case ens_gateway_interfaces.QueryKindText:
		result, err = u.resolveText(ctx, query)
	case ens_gateway_interfaces.QueryKindAddr:
		result, err = u.resolveAddr(ctx, query)
	default:
		return nil, fmt.Errorf("unsupported query kind: %d", query.Kind)
	}
	if err != nil {
		return nil, err
	}

	signed, err := u.responseBuilder.Build(ctx, request.Sender, request.Data, result)
	if err != nil {
		return nil, fmt.Errorf("build signed response: %w", err)
	}

	return &responses.CcipReadResponse{Data: hexutil.Encode(signed)}, nil
}

func (u *ResolveOffchainQuery) resolveText(ctx context.Context, query ens_gateway_interfaces.DecodedQuery) ([]byte, error) {
	value, err := u.textValue(ctx, query)
	if err != nil {
		return nil, err
	}

	encoded, err := u.resolverServiceCodec.EncodeTextResult(value)
	if err != nil {
		return nil, fmt.Errorf("encode text result: %w", err)
	}
	return encoded, nil
}

func (u *ResolveOffchainQuery) textValue(ctx context.Context, query ens_gateway_interfaces.DecodedQuery) (string, error) {
	if cached := u.resolutionCache.GetText(query.ChainId, query.Node, query.Key); cached != nil {
		return *cached, nil
	}

	value, err := u.agentEnsRegistryReader.Text(ctx, query.ChainId, query.Node, query.Key)
	if err != nil {
		return "", fmt.Errorf("read agent ens registry: %w", err)
	}
	if value == "" {
		fallback, fallbackErr := u.universalResolverClient.ResolveText(ctx, query.ChainId, query.DnsName, query.InnerData)
		switch {
		case fallbackErr != nil:
			u.logger.Warn("L1 universal resolver fallback failed",
				zap.Error(fallbackErr),
				zap.Uint64("chain_id", query.ChainId),
				zap.String("key", query.Key),
			)
		case fallback != "":
			value = fallback
		}
	}

	u.resolutionCache.SetText(query.ChainId, query.Node, query.Key, value)
	return value, nil
}

func (u *ResolveOffchainQuery) resolveAddr(ctx context.Context, query ens_gateway_interfaces.DecodedQuery) ([]byte, error) {
	value, err := u.addrValue(ctx, query)
	if err != nil {
		return nil, err
	}

	encoded, err := u.resolverServiceCodec.EncodeAddrResult(value)
	if err != nil {
		return nil, fmt.Errorf("encode addr result: %w", err)
	}
	return encoded, nil
}

func (u *ResolveOffchainQuery) addrValue(ctx context.Context, query ens_gateway_interfaces.DecodedQuery) (common.Address, error) {
	if cached := u.resolutionCache.GetAddr(query.ChainId, query.Node); cached != nil {
		return *cached, nil
	}

	value, err := u.agentEnsRegistryReader.Addr(ctx, query.ChainId, query.Node)
	if err != nil {
		return common.Address{}, fmt.Errorf("read agent ens registry: %w", err)
	}
	if value == (common.Address{}) {
		fallback, fallbackErr := u.universalResolverClient.ResolveAddr(ctx, query.ChainId, query.DnsName, query.InnerData)
		switch {
		case fallbackErr != nil:
			u.logger.Warn("L1 universal resolver fallback failed",
				zap.Error(fallbackErr),
				zap.Uint64("chain_id", query.ChainId),
			)
		case fallback != (common.Address{}):
			value = fallback
		}
	}

	u.resolutionCache.SetAddr(query.ChainId, query.Node, value)
	return value, nil
}

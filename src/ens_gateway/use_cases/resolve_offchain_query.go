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
	logger            *zap.Logger
	codec             ens_gateway_interfaces.ResolverServiceCodecInterface
	reader            ens_gateway_interfaces.AgentEnsRegistryReaderInterface
	universalResolver ens_gateway_interfaces.UniversalResolverClientInterface
	builder           ens_gateway_interfaces.ResponseBuilderInterface
}

func NewResolveOffchainQuery(
	logger *zap.Logger,
	codec ens_gateway_interfaces.ResolverServiceCodecInterface,
	reader ens_gateway_interfaces.AgentEnsRegistryReaderInterface,
	universalResolver ens_gateway_interfaces.UniversalResolverClientInterface,
	builder ens_gateway_interfaces.ResponseBuilderInterface,
) *ResolveOffchainQuery {
	return &ResolveOffchainQuery{
		logger:            logger,
		codec:             codec,
		reader:            reader,
		universalResolver: universalResolver,
		builder:           builder,
	}
}

func (u *ResolveOffchainQuery) Execute(
	ctx context.Context,
	request *requests.ResolveOffchainQueryRequest,
) (*responses.CcipReadResponse, error) {
	query, err := u.codec.Decode(request.Data)
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

	signed, err := u.builder.Build(ctx, request.Sender, request.Data, result)
	if err != nil {
		return nil, fmt.Errorf("build signed response: %w", err)
	}

	return &responses.CcipReadResponse{Data: hexutil.Encode(signed)}, nil
}

func (u *ResolveOffchainQuery) resolveText(ctx context.Context, query ens_gateway_interfaces.DecodedQuery) ([]byte, error) {
	value, err := u.reader.Text(ctx, query.ChainId, query.Node, query.Key)
	if err != nil {
		return nil, fmt.Errorf("read agent ens registry: %w", err)
	}
	if value == "" {
		fallback, fallbackErr := u.universalResolver.ResolveText(ctx, query.ChainId, query.DnsName, query.InnerData)
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

	encoded, err := u.codec.EncodeTextResult(value)
	if err != nil {
		return nil, fmt.Errorf("encode text result: %w", err)
	}
	return encoded, nil
}

func (u *ResolveOffchainQuery) resolveAddr(ctx context.Context, query ens_gateway_interfaces.DecodedQuery) ([]byte, error) {
	value, err := u.reader.Addr(ctx, query.ChainId, query.Node)
	if err != nil {
		return nil, fmt.Errorf("read agent ens registry: %w", err)
	}
	if value == (common.Address{}) {
		fallback, fallbackErr := u.universalResolver.ResolveAddr(ctx, query.ChainId, query.DnsName, query.InnerData)
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

	encoded, err := u.codec.EncodeAddrResult(value)
	if err != nil {
		return nil, fmt.Errorf("encode addr result: %w", err)
	}
	return encoded, nil
}

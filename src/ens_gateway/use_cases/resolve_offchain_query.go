package use_cases

import (
	"context"
	"fmt"

	ens_gateway_interfaces "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/services/interfaces"
	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/use_cases/responses"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type ResolveOffchainQueryInterface interface {
	Execute(ctx context.Context, request *requests.ResolveOffchainQueryRequest) (*responses.CcipReadResponse, error)
}

// ResolveOffchainQuery is the CCIP-Read vertical slice: decode the resolver
// call data, read the record from AgentEnsRegistry, and return a signed
// response the L1 resolver can authenticate.
type ResolveOffchainQuery struct {
	codec   ens_gateway_interfaces.ResolverServiceCodecInterface
	reader  ens_gateway_interfaces.AgentEnsRegistryReaderInterface
	builder ens_gateway_interfaces.ResponseBuilderInterface
}

func NewResolveOffchainQuery(
	codec ens_gateway_interfaces.ResolverServiceCodecInterface,
	reader ens_gateway_interfaces.AgentEnsRegistryReaderInterface,
	builder ens_gateway_interfaces.ResponseBuilderInterface,
) *ResolveOffchainQuery {
	return &ResolveOffchainQuery{codec: codec, reader: reader, builder: builder}
}

func (u *ResolveOffchainQuery) Execute(
	ctx context.Context,
	request *requests.ResolveOffchainQueryRequest,
) (*responses.CcipReadResponse, error) {
	chainId, node, key, err := u.codec.Decode(request.Data)
	if err != nil {
		return nil, fmt.Errorf("decode call data: %w", err)
	}

	value, err := u.reader.Text(ctx, chainId, node, key)
	if err != nil {
		return nil, fmt.Errorf("read agent ens registry: %w", err)
	}

	result, err := u.codec.EncodeResult(value)
	if err != nil {
		return nil, fmt.Errorf("encode result: %w", err)
	}

	signed, err := u.builder.Build(ctx, request.Sender, request.Data, result)
	if err != nil {
		return nil, fmt.Errorf("build signed response: %w", err)
	}

	return &responses.CcipReadResponse{Data: hexutil.Encode(signed)}, nil
}

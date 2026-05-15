package use_cases_test

import (
	"context"
	"errors"
	"testing"

	gateway_use_cases "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/use_cases"
	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/use_cases/requests"
	interfaces_mocks "github.com/6022-labs/agentic-ens-gateway/tests/ens_gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenResolvingOffchainQueryTestingSuite struct {
	mockCodec   *interfaces_mocks.MockResolverServiceCodecInterface
	mockReader  *interfaces_mocks.MockAgentEnsRegistryReaderInterface
	mockBuilder *interfaces_mocks.MockResponseBuilderInterface
	sut         *gateway_use_cases.ResolveOffchainQuery
}

func WhenResolvingOffchainQueryBeforeEach(t *testing.T) *WhenResolvingOffchainQueryTestingSuite {
	ctrl := gomock.NewController(t)
	mockCodec := interfaces_mocks.NewMockResolverServiceCodecInterface(ctrl)
	mockReader := interfaces_mocks.NewMockAgentEnsRegistryReaderInterface(ctrl)
	mockBuilder := interfaces_mocks.NewMockResponseBuilderInterface(ctrl)

	return &WhenResolvingOffchainQueryTestingSuite{
		mockCodec:   mockCodec,
		mockReader:  mockReader,
		mockBuilder: mockBuilder,
		sut:         gateway_use_cases.NewResolveOffchainQuery(mockCodec, mockReader, mockBuilder),
	}
}

func TestWhenResolvingOffchainQuery(t *testing.T) {
	t.Parallel()

	request := &requests.ResolveOffchainQueryRequest{
		Sender: common.HexToAddress("0x1111111111111111111111111111111111111111"),
		Data:   []byte{0x01, 0x02, 0x03},
	}
	var node [32]byte
	node[31] = 0x07
	const (
		chainId uint64 = 137
		key            = "agent.runtime-url"
	)

	t.Run("Given the codec, reader and builder all succeed", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingOffchainQueryTestingSuite) {
			suite.mockCodec.EXPECT().Decode(request.Data).Return(chainId, node, key, nil)
			suite.mockReader.EXPECT().Text(gomock.Any(), chainId, node, key).Return("https://agent.example/", nil)
			suite.mockCodec.EXPECT().EncodeResult("https://agent.example/").Return([]byte{0xde, 0xad}, nil)
			suite.mockBuilder.EXPECT().
				Build(gomock.Any(), request.Sender, request.Data, []byte{0xde, 0xad}).
				Return([]byte{0xbe, 0xef}, nil)
		}

		t.Run("Should return the builder's signed response as hex", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingOffchainQueryBeforeEach(t)
			initSuite(suite)

			response, err := suite.sut.Execute(context.Background(), request)

			assert.NoError(t, err)
			assert.Equal(t, "0xbeef", response.Data)
		})
	})

	t.Run("Given the codec fails to decode", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingOffchainQueryTestingSuite) {
			suite.mockCodec.EXPECT().Decode(request.Data).Return(uint64(0), [32]byte{}, "", errors.New("bad selector"))
		}

		t.Run("Should return the error without reading the registry", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingOffchainQueryBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Execute(context.Background(), request)

			assert.Error(t, err)
		})
	})

	t.Run("Given the registry read fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingOffchainQueryTestingSuite) {
			suite.mockCodec.EXPECT().Decode(request.Data).Return(chainId, node, key, nil)
			suite.mockReader.EXPECT().Text(gomock.Any(), chainId, node, key).Return("", errors.New("rpc down"))
		}

		t.Run("Should return the error without building a response", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingOffchainQueryBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Execute(context.Background(), request)

			assert.Error(t, err)
		})
	})

	t.Run("Given the codec fails to encode the result", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingOffchainQueryTestingSuite) {
			suite.mockCodec.EXPECT().Decode(request.Data).Return(chainId, node, key, nil)
			suite.mockReader.EXPECT().Text(gomock.Any(), chainId, node, key).Return("https://agent.example/", nil)
			suite.mockCodec.EXPECT().EncodeResult("https://agent.example/").Return(nil, errors.New("encode failed"))
		}

		t.Run("Should return the error without building a response", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingOffchainQueryBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Execute(context.Background(), request)

			assert.Error(t, err)
		})
	})

	t.Run("Given the response builder fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingOffchainQueryTestingSuite) {
			suite.mockCodec.EXPECT().Decode(request.Data).Return(chainId, node, key, nil)
			suite.mockReader.EXPECT().Text(gomock.Any(), chainId, node, key).Return("https://agent.example/", nil)
			suite.mockCodec.EXPECT().EncodeResult("https://agent.example/").Return([]byte{0xde, 0xad}, nil)
			suite.mockBuilder.EXPECT().
				Build(gomock.Any(), request.Sender, request.Data, []byte{0xde, 0xad}).
				Return(nil, errors.New("sign failed"))
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingOffchainQueryBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Execute(context.Background(), request)

			assert.Error(t, err)
		})
	})
}

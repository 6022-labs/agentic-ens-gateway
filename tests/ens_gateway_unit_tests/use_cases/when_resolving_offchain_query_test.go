package use_cases_test

import (
	"context"
	"errors"
	"testing"

	ens_gateway_interfaces "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/services/interfaces"
	gateway_use_cases "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/use_cases"
	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/use_cases/requests"
	interfaces_mocks "github.com/6022-labs/agentic-ens-gateway/tests/ens_gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type WhenResolvingOffchainQueryTestingSuite struct {
	mockCodec             *interfaces_mocks.MockResolverServiceCodecInterface
	mockReader            *interfaces_mocks.MockAgentEnsRegistryReaderInterface
	mockUniversalResolver *interfaces_mocks.MockUniversalResolverClientInterface
	mockBuilder           *interfaces_mocks.MockResponseBuilderInterface
	mockCache             *interfaces_mocks.MockResolutionCacheInterface
	logs                  *observer.ObservedLogs
	sut                   *gateway_use_cases.ResolveOffchainQuery
}

func WhenResolvingOffchainQueryBeforeEach(t *testing.T) *WhenResolvingOffchainQueryTestingSuite {
	ctrl := gomock.NewController(t)
	mockCodec := interfaces_mocks.NewMockResolverServiceCodecInterface(ctrl)
	mockReader := interfaces_mocks.NewMockAgentEnsRegistryReaderInterface(ctrl)
	mockUniversalResolver := interfaces_mocks.NewMockUniversalResolverClientInterface(ctrl)
	mockBuilder := interfaces_mocks.NewMockResponseBuilderInterface(ctrl)
	mockCache := interfaces_mocks.NewMockResolutionCacheInterface(ctrl)

	core, logs := observer.New(zapcore.WarnLevel)
	logger := zap.New(core)

	return &WhenResolvingOffchainQueryTestingSuite{
		mockCodec:             mockCodec,
		mockReader:            mockReader,
		mockUniversalResolver: mockUniversalResolver,
		mockBuilder:           mockBuilder,
		mockCache:             mockCache,
		logs:                  logs,
		sut:                   gateway_use_cases.NewResolveOffchainQuery(logger, mockCodec, mockReader, mockUniversalResolver, mockBuilder, mockCache),
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
		key            = "url"
	)
	dnsName := []byte{0x04, 'a', 'g', 'n', 't', 0x03, '1', '3', '7', 0x00}
	innerData := []byte{0xaa, 0xbb, 0xcc}
	addrInnerData := []byte{0xdd, 0xee, 0xff}

	decoded := ens_gateway_interfaces.DecodedQuery{
		Kind:      ens_gateway_interfaces.QueryKindText,
		ChainId:   chainId,
		Node:      node,
		Key:       key,
		DnsName:   dnsName,
		InnerData: innerData,
	}

	decodedAddr := ens_gateway_interfaces.DecodedQuery{
		Kind:      ens_gateway_interfaces.QueryKindAddr,
		ChainId:   chainId,
		Node:      node,
		DnsName:   dnsName,
		InnerData: addrInnerData,
	}

	l2Addr := common.HexToAddress("0x2222222222222222222222222222222222222222")
	l1Addr := common.HexToAddress("0x3333333333333333333333333333333333333333")

	t.Run("Given the text value is cached", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingOffchainQueryTestingSuite) {
			cached := "https://cached.example/"
			suite.mockCodec.EXPECT().Decode(request.Data).Return(decoded, nil)
			suite.mockCache.EXPECT().GetText(chainId, node, key).Return(&cached)
			suite.mockCodec.EXPECT().EncodeTextResult(cached).Return([]byte{0xde, 0xad}, nil)
			suite.mockBuilder.EXPECT().
				Build(gomock.Any(), request.Sender, request.Data, []byte{0xde, 0xad}).
				Return([]byte{0xbe, 0xef}, nil)
		}

		t.Run("Should return the cached value without reading the registry or the fallback", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingOffchainQueryBeforeEach(t)
			initSuite(suite)

			response, err := suite.sut.Execute(context.Background(), request)

			assert.NoError(t, err)
			assert.Equal(t, "0xbeef", response.Data)
		})
	})

	t.Run("Given the L2 registry has the record", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingOffchainQueryTestingSuite) {
			suite.mockCodec.EXPECT().Decode(request.Data).Return(decoded, nil)
			suite.mockCache.EXPECT().GetText(chainId, node, key).Return(nil)
			suite.mockReader.EXPECT().Text(gomock.Any(), chainId, node, key).Return("https://agent.example/", nil)
			suite.mockCache.EXPECT().SetText(chainId, node, key, "https://agent.example/")
			suite.mockCodec.EXPECT().EncodeTextResult("https://agent.example/").Return([]byte{0xde, 0xad}, nil)
			suite.mockBuilder.EXPECT().
				Build(gomock.Any(), request.Sender, request.Data, []byte{0xde, 0xad}).
				Return([]byte{0xbe, 0xef}, nil)
		}

		t.Run("Should cache the value, sign it, and not call the L1 fallback", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingOffchainQueryBeforeEach(t)
			initSuite(suite)

			response, err := suite.sut.Execute(context.Background(), request)

			assert.NoError(t, err)
			assert.Equal(t, "0xbeef", response.Data)
		})
	})

	t.Run("Given the L2 registry returns empty and the L1 fallback returns a value", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingOffchainQueryTestingSuite) {
			suite.mockCodec.EXPECT().Decode(request.Data).Return(decoded, nil)
			suite.mockCache.EXPECT().GetText(chainId, node, key).Return(nil)
			suite.mockReader.EXPECT().Text(gomock.Any(), chainId, node, key).Return("", nil)
			suite.mockUniversalResolver.EXPECT().ResolveText(gomock.Any(), chainId, dnsName, innerData).Return("https://l1.example/", nil)
			suite.mockCache.EXPECT().SetText(chainId, node, key, "https://l1.example/")
			suite.mockCodec.EXPECT().EncodeTextResult("https://l1.example/").Return([]byte{0xab, 0xcd}, nil)
			suite.mockBuilder.EXPECT().
				Build(gomock.Any(), request.Sender, request.Data, []byte{0xab, 0xcd}).
				Return([]byte{0xfa, 0xce}, nil)
		}

		t.Run("Should sign the L1 fallback value", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingOffchainQueryBeforeEach(t)
			initSuite(suite)

			response, err := suite.sut.Execute(context.Background(), request)

			assert.NoError(t, err)
			assert.Equal(t, "0xface", response.Data)
		})
	})

	t.Run("Given both the L2 registry and the L1 fallback return empty", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingOffchainQueryTestingSuite) {
			suite.mockCodec.EXPECT().Decode(request.Data).Return(decoded, nil)
			suite.mockCache.EXPECT().GetText(chainId, node, key).Return(nil)
			suite.mockReader.EXPECT().Text(gomock.Any(), chainId, node, key).Return("", nil)
			suite.mockUniversalResolver.EXPECT().ResolveText(gomock.Any(), chainId, dnsName, innerData).Return("", nil)
			suite.mockCache.EXPECT().SetText(chainId, node, key, "")
			suite.mockCodec.EXPECT().EncodeTextResult("").Return([]byte{0x00}, nil)
			suite.mockBuilder.EXPECT().
				Build(gomock.Any(), request.Sender, request.Data, []byte{0x00}).
				Return([]byte{0x42}, nil)
		}

		t.Run("Should sign an empty value without erroring", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingOffchainQueryBeforeEach(t)
			initSuite(suite)

			response, err := suite.sut.Execute(context.Background(), request)

			assert.NoError(t, err)
			assert.Equal(t, "0x42", response.Data)
		})
	})

	t.Run("Given the L2 registry returns empty and the L1 fallback errors", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingOffchainQueryTestingSuite) {
			suite.mockCodec.EXPECT().Decode(request.Data).Return(decoded, nil)
			suite.mockCache.EXPECT().GetText(chainId, node, key).Return(nil)
			suite.mockReader.EXPECT().Text(gomock.Any(), chainId, node, key).Return("", nil)
			suite.mockUniversalResolver.EXPECT().ResolveText(gomock.Any(), chainId, dnsName, innerData).Return("", errors.New("OffchainLookup"))
			suite.mockCache.EXPECT().SetText(chainId, node, key, "")
			suite.mockCodec.EXPECT().EncodeTextResult("").Return([]byte{0x00}, nil)
			suite.mockBuilder.EXPECT().
				Build(gomock.Any(), request.Sender, request.Data, []byte{0x00}).
				Return([]byte{0x42}, nil)
		}

		t.Run("Should swallow the error, log a warning, and sign an empty value", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingOffchainQueryBeforeEach(t)
			initSuite(suite)

			response, err := suite.sut.Execute(context.Background(), request)

			assert.NoError(t, err)
			assert.Equal(t, "0x42", response.Data)
			warns := suite.logs.FilterMessage("L1 universal resolver fallback failed").All()
			assert.Len(t, warns, 1)
		})
	})

	t.Run("Given the codec fails to decode", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingOffchainQueryTestingSuite) {
			suite.mockCodec.EXPECT().Decode(request.Data).Return(ens_gateway_interfaces.DecodedQuery{}, errors.New("bad selector"))
		}

		t.Run("Should return the error without touching the cache, registry, or fallback", func(t *testing.T) {
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
			suite.mockCodec.EXPECT().Decode(request.Data).Return(decoded, nil)
			suite.mockCache.EXPECT().GetText(chainId, node, key).Return(nil)
			suite.mockReader.EXPECT().Text(gomock.Any(), chainId, node, key).Return("", errors.New("rpc down"))
		}

		t.Run("Should return the error without caching, invoking the fallback, or building a response", func(t *testing.T) {
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
			suite.mockCodec.EXPECT().Decode(request.Data).Return(decoded, nil)
			suite.mockCache.EXPECT().GetText(chainId, node, key).Return(nil)
			suite.mockReader.EXPECT().Text(gomock.Any(), chainId, node, key).Return("https://agent.example/", nil)
			suite.mockCache.EXPECT().SetText(chainId, node, key, "https://agent.example/")
			suite.mockCodec.EXPECT().EncodeTextResult("https://agent.example/").Return(nil, errors.New("encode failed"))
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
			suite.mockCodec.EXPECT().Decode(request.Data).Return(decoded, nil)
			suite.mockCache.EXPECT().GetText(chainId, node, key).Return(nil)
			suite.mockReader.EXPECT().Text(gomock.Any(), chainId, node, key).Return("https://agent.example/", nil)
			suite.mockCache.EXPECT().SetText(chainId, node, key, "https://agent.example/")
			suite.mockCodec.EXPECT().EncodeTextResult("https://agent.example/").Return([]byte{0xde, 0xad}, nil)
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

	t.Run("Given the addr value is cached", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingOffchainQueryTestingSuite) {
			suite.mockCodec.EXPECT().Decode(request.Data).Return(decodedAddr, nil)
			suite.mockCache.EXPECT().GetAddr(chainId, node).Return(&l2Addr)
			suite.mockCodec.EXPECT().EncodeAddrResult(l2Addr).Return([]byte{0xa1, 0xa2}, nil)
			suite.mockBuilder.EXPECT().
				Build(gomock.Any(), request.Sender, request.Data, []byte{0xa1, 0xa2}).
				Return([]byte{0xb1, 0xb2}, nil)
		}

		t.Run("Should return the cached address without reading the registry or the fallback", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingOffchainQueryBeforeEach(t)
			initSuite(suite)

			response, err := suite.sut.Execute(context.Background(), request)

			assert.NoError(t, err)
			assert.Equal(t, "0xb1b2", response.Data)
		})
	})

	t.Run("Given an addr query and the L2 registry has the address", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingOffchainQueryTestingSuite) {
			suite.mockCodec.EXPECT().Decode(request.Data).Return(decodedAddr, nil)
			suite.mockCache.EXPECT().GetAddr(chainId, node).Return(nil)
			suite.mockReader.EXPECT().Addr(gomock.Any(), chainId, node).Return(l2Addr, nil)
			suite.mockCache.EXPECT().SetAddr(chainId, node, l2Addr)
			suite.mockCodec.EXPECT().EncodeAddrResult(l2Addr).Return([]byte{0xa1, 0xa2}, nil)
			suite.mockBuilder.EXPECT().
				Build(gomock.Any(), request.Sender, request.Data, []byte{0xa1, 0xa2}).
				Return([]byte{0xb1, 0xb2}, nil)
		}

		t.Run("Should cache and sign the L2 address and not call the L1 fallback", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingOffchainQueryBeforeEach(t)
			initSuite(suite)

			response, err := suite.sut.Execute(context.Background(), request)

			assert.NoError(t, err)
			assert.Equal(t, "0xb1b2", response.Data)
		})
	})

	t.Run("Given an addr query and the L2 registry returns the zero address but the L1 fallback resolves", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingOffchainQueryTestingSuite) {
			suite.mockCodec.EXPECT().Decode(request.Data).Return(decodedAddr, nil)
			suite.mockCache.EXPECT().GetAddr(chainId, node).Return(nil)
			suite.mockReader.EXPECT().Addr(gomock.Any(), chainId, node).Return(common.Address{}, nil)
			suite.mockUniversalResolver.EXPECT().ResolveAddr(gomock.Any(), chainId, dnsName, addrInnerData).Return(l1Addr, nil)
			suite.mockCache.EXPECT().SetAddr(chainId, node, l1Addr)
			suite.mockCodec.EXPECT().EncodeAddrResult(l1Addr).Return([]byte{0xc1, 0xc2}, nil)
			suite.mockBuilder.EXPECT().
				Build(gomock.Any(), request.Sender, request.Data, []byte{0xc1, 0xc2}).
				Return([]byte{0xd1, 0xd2}, nil)
		}

		t.Run("Should sign the L1 fallback address", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingOffchainQueryBeforeEach(t)
			initSuite(suite)

			response, err := suite.sut.Execute(context.Background(), request)

			assert.NoError(t, err)
			assert.Equal(t, "0xd1d2", response.Data)
		})
	})

	t.Run("Given an addr query and both L2 and L1 return the zero address", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingOffchainQueryTestingSuite) {
			suite.mockCodec.EXPECT().Decode(request.Data).Return(decodedAddr, nil)
			suite.mockCache.EXPECT().GetAddr(chainId, node).Return(nil)
			suite.mockReader.EXPECT().Addr(gomock.Any(), chainId, node).Return(common.Address{}, nil)
			suite.mockUniversalResolver.EXPECT().ResolveAddr(gomock.Any(), chainId, dnsName, addrInnerData).Return(common.Address{}, nil)
			suite.mockCache.EXPECT().SetAddr(chainId, node, common.Address{})
			suite.mockCodec.EXPECT().EncodeAddrResult(common.Address{}).Return([]byte{0x00}, nil)
			suite.mockBuilder.EXPECT().
				Build(gomock.Any(), request.Sender, request.Data, []byte{0x00}).
				Return([]byte{0x42}, nil)
		}

		t.Run("Should sign the zero address without erroring", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingOffchainQueryBeforeEach(t)
			initSuite(suite)

			response, err := suite.sut.Execute(context.Background(), request)

			assert.NoError(t, err)
			assert.Equal(t, "0x42", response.Data)
		})
	})

	t.Run("Given an addr query and the L1 fallback errors", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingOffchainQueryTestingSuite) {
			suite.mockCodec.EXPECT().Decode(request.Data).Return(decodedAddr, nil)
			suite.mockCache.EXPECT().GetAddr(chainId, node).Return(nil)
			suite.mockReader.EXPECT().Addr(gomock.Any(), chainId, node).Return(common.Address{}, nil)
			suite.mockUniversalResolver.EXPECT().ResolveAddr(gomock.Any(), chainId, dnsName, addrInnerData).Return(common.Address{}, errors.New("OffchainLookup"))
			suite.mockCache.EXPECT().SetAddr(chainId, node, common.Address{})
			suite.mockCodec.EXPECT().EncodeAddrResult(common.Address{}).Return([]byte{0x00}, nil)
			suite.mockBuilder.EXPECT().
				Build(gomock.Any(), request.Sender, request.Data, []byte{0x00}).
				Return([]byte{0x42}, nil)
		}

		t.Run("Should swallow the error, log a warning, and sign the zero address", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingOffchainQueryBeforeEach(t)
			initSuite(suite)

			response, err := suite.sut.Execute(context.Background(), request)

			assert.NoError(t, err)
			assert.Equal(t, "0x42", response.Data)
			warns := suite.logs.FilterMessage("L1 universal resolver fallback failed").All()
			assert.Len(t, warns, 1)
		})
	})

	t.Run("Given an addr query and the registry read fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingOffchainQueryTestingSuite) {
			suite.mockCodec.EXPECT().Decode(request.Data).Return(decodedAddr, nil)
			suite.mockCache.EXPECT().GetAddr(chainId, node).Return(nil)
			suite.mockReader.EXPECT().Addr(gomock.Any(), chainId, node).Return(common.Address{}, errors.New("rpc down"))
		}

		t.Run("Should return the error without caching, invoking the fallback, or building a response", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingOffchainQueryBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Execute(context.Background(), request)

			assert.Error(t, err)
		})
	})
}

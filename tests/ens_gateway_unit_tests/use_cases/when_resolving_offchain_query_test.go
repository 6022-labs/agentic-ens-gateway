package use_cases_test

import (
	"context"
	"errors"
	"testing"

	apperrors "github.com/6022-labs/agentic-ens-gateway/src/common/errors"
	metrics_interfaces "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/metrics/interfaces"
	ens_gateway_interfaces "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/services/interfaces"
	traces_interfaces "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/traces/interfaces"
	gateway_use_cases "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/use_cases"
	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/use_cases/requests"
	metrics_mocks "github.com/6022-labs/agentic-ens-gateway/tests/ens_gateway_mocks/metrics_mocks/interfaces_mocks"
	interfaces_mocks "github.com/6022-labs/agentic-ens-gateway/tests/ens_gateway_mocks/services_mocks/interfaces_mocks"
	traces_mocks "github.com/6022-labs/agentic-ens-gateway/tests/ens_gateway_mocks/traces_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type WhenResolvingOffchainQueryTestingSuite struct {
	mockResolverServiceCodec    *interfaces_mocks.MockResolverServiceCodecInterface
	mockAgentEnsRegistryReader  *interfaces_mocks.MockAgentEnsRegistryReaderInterface
	mockUniversalResolverClient *interfaces_mocks.MockUniversalResolverClientInterface
	mockResponseBuilder         *interfaces_mocks.MockResponseBuilderInterface
	mockResolutionCache         *interfaces_mocks.MockResolutionCacheInterface
	mockEnsResolutionMetrics    *metrics_mocks.MockEnsResolutionMetricsInterface
	logs                        *observer.ObservedLogs
	sut                         *gateway_use_cases.ResolveOffchainQuery
}

func WhenResolvingOffchainQueryBeforeEach(t *testing.T) *WhenResolvingOffchainQueryTestingSuite {
	ctrl := gomock.NewController(t)
	mockResolverServiceCodec := interfaces_mocks.NewMockResolverServiceCodecInterface(ctrl)
	mockAgentEnsRegistryReader := interfaces_mocks.NewMockAgentEnsRegistryReaderInterface(ctrl)
	mockUniversalResolverClient := interfaces_mocks.NewMockUniversalResolverClientInterface(ctrl)
	mockResponseBuilder := interfaces_mocks.NewMockResponseBuilderInterface(ctrl)
	mockResolutionCache := interfaces_mocks.NewMockResolutionCacheInterface(ctrl)

	mockEnsResolutionMetrics := metrics_mocks.NewMockEnsResolutionMetricsInterface(ctrl)

	mockSpan := traces_mocks.NewMockSpan(ctrl)
	mockSpan.EXPECT().Fail(gomock.Any()).AnyTimes()
	mockSpan.EXPECT().End().AnyTimes()

	mockResolutionSpan := traces_mocks.NewMockResolutionSpan(ctrl)
	mockResolutionSpan.EXPECT().RecordSource(gomock.Any()).AnyTimes()
	mockResolutionSpan.EXPECT().Fail(gomock.Any()).AnyTimes()
	mockResolutionSpan.EXPECT().End().AnyTimes()

	mockResolutionTracer := traces_mocks.NewMockResolutionTracerInterface(ctrl)
	mockResolutionTracer.EXPECT().StartResolution(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, _ string, _ uint64) (context.Context, traces_interfaces.ResolutionSpan) {
			return ctx, mockResolutionSpan
		}).AnyTimes()
	mockResolutionTracer.EXPECT().StartRegistryRead(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, _ uint64) (context.Context, traces_interfaces.Span) {
			return ctx, mockSpan
		}).AnyTimes()
	mockResolutionTracer.EXPECT().StartUniversalResolverFallback(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, _ uint64) (context.Context, traces_interfaces.Span) {
			return ctx, mockSpan
		}).AnyTimes()

	core, logs := observer.New(zapcore.WarnLevel)
	logger := zap.New(core)

	return &WhenResolvingOffchainQueryTestingSuite{
		mockResolverServiceCodec:    mockResolverServiceCodec,
		mockAgentEnsRegistryReader:  mockAgentEnsRegistryReader,
		mockUniversalResolverClient: mockUniversalResolverClient,
		mockResponseBuilder:         mockResponseBuilder,
		mockResolutionCache:         mockResolutionCache,
		mockEnsResolutionMetrics:    mockEnsResolutionMetrics,
		logs:                        logs,
		sut: gateway_use_cases.NewResolveOffchainQuery(
			logger,
			mockResolverServiceCodec,
			mockAgentEnsRegistryReader,
			mockUniversalResolverClient,
			mockResponseBuilder,
			mockResolutionCache,
			mockEnsResolutionMetrics,
			mockResolutionTracer,
		),
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
			suite.mockResolverServiceCodec.EXPECT().Decode(request.Data).Return(decoded, nil)
			suite.mockResolutionCache.EXPECT().GetText(chainId, node, key).Return(&cached)
			suite.mockResolverServiceCodec.EXPECT().EncodeTextResult(cached).Return([]byte{0xde, 0xad}, nil)
			suite.mockResponseBuilder.EXPECT().
				Build(gomock.Any(), request.Sender, request.Data, []byte{0xde, 0xad}).
				Return([]byte{0xbe, 0xef}, nil)

			suite.mockEnsResolutionMetrics.EXPECT().RecordCacheResult(gomock.Any(), "text", true)
			suite.mockEnsResolutionMetrics.EXPECT().RecordResolution(
				gomock.Any(), "text", metrics_interfaces.ResolutionSourceCache, chainId, gomock.Any(), nil)
			suite.mockEnsResolutionMetrics.EXPECT().RecordSignedResponse(gomock.Any(), nil)
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
			suite.mockResolverServiceCodec.EXPECT().Decode(request.Data).Return(decoded, nil)
			suite.mockResolutionCache.EXPECT().GetText(chainId, node, key).Return(nil)
			suite.mockAgentEnsRegistryReader.EXPECT().
				Text(gomock.Any(), chainId, node, key).
				Return("https://agent.example/", nil)
			suite.mockResolutionCache.EXPECT().SetText(chainId, node, key, "https://agent.example/")
			suite.mockResolverServiceCodec.EXPECT().
				EncodeTextResult("https://agent.example/").
				Return([]byte{0xde, 0xad}, nil)
			suite.mockResponseBuilder.EXPECT().
				Build(gomock.Any(), request.Sender, request.Data, []byte{0xde, 0xad}).
				Return([]byte{0xbe, 0xef}, nil)

			suite.mockEnsResolutionMetrics.EXPECT().RecordCacheResult(gomock.Any(), "text", false)
			suite.mockEnsResolutionMetrics.EXPECT().RecordResolution(
				gomock.Any(), "text", metrics_interfaces.ResolutionSourceAgentEnsRegistry, chainId, gomock.Any(), nil)
			suite.mockEnsResolutionMetrics.EXPECT().RecordSignedResponse(gomock.Any(), nil)
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
			suite.mockResolverServiceCodec.EXPECT().Decode(request.Data).Return(decoded, nil)
			suite.mockResolutionCache.EXPECT().GetText(chainId, node, key).Return(nil)
			suite.mockAgentEnsRegistryReader.EXPECT().Text(gomock.Any(), chainId, node, key).Return("", nil)
			suite.mockUniversalResolverClient.EXPECT().
				ResolveText(gomock.Any(), chainId, dnsName, innerData).
				Return("https://l1.example/", nil)
			suite.mockResolutionCache.EXPECT().SetText(chainId, node, key, "https://l1.example/")
			suite.mockResolverServiceCodec.EXPECT().
				EncodeTextResult("https://l1.example/").
				Return([]byte{0xab, 0xcd}, nil)
			suite.mockResponseBuilder.EXPECT().
				Build(gomock.Any(), request.Sender, request.Data, []byte{0xab, 0xcd}).
				Return([]byte{0xfa, 0xce}, nil)

			suite.mockEnsResolutionMetrics.EXPECT().RecordCacheResult(gomock.Any(), "text", false)
			suite.mockEnsResolutionMetrics.EXPECT().RecordFallback(
				gomock.Any(), "text", metrics_interfaces.FallbackOutcomeResolved, chainId)
			suite.mockEnsResolutionMetrics.EXPECT().RecordResolution(
				gomock.Any(), "text", metrics_interfaces.ResolutionSourceUniversalResolver, chainId, gomock.Any(), nil)
			suite.mockEnsResolutionMetrics.EXPECT().RecordSignedResponse(gomock.Any(), nil)
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
			suite.mockResolverServiceCodec.EXPECT().Decode(request.Data).Return(decoded, nil)
			suite.mockResolutionCache.EXPECT().GetText(chainId, node, key).Return(nil)
			suite.mockAgentEnsRegistryReader.EXPECT().Text(gomock.Any(), chainId, node, key).Return("", nil)
			suite.mockUniversalResolverClient.EXPECT().
				ResolveText(gomock.Any(), chainId, dnsName, innerData).
				Return("", nil)
			suite.mockResolutionCache.EXPECT().SetText(chainId, node, key, "")
			suite.mockResolverServiceCodec.EXPECT().EncodeTextResult("").Return([]byte{0x00}, nil)
			suite.mockResponseBuilder.EXPECT().
				Build(gomock.Any(), request.Sender, request.Data, []byte{0x00}).
				Return([]byte{0x42}, nil)

			suite.mockEnsResolutionMetrics.EXPECT().RecordCacheResult(gomock.Any(), "text", false)
			suite.mockEnsResolutionMetrics.EXPECT().RecordFallback(
				gomock.Any(), "text", metrics_interfaces.FallbackOutcomeEmpty, chainId)
			suite.mockEnsResolutionMetrics.EXPECT().RecordResolution(
				gomock.Any(), "text", metrics_interfaces.ResolutionSourceAgentEnsRegistry, chainId, gomock.Any(), nil)
			suite.mockEnsResolutionMetrics.EXPECT().RecordSignedResponse(gomock.Any(), nil)
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
			suite.mockResolverServiceCodec.EXPECT().Decode(request.Data).Return(decoded, nil)
			suite.mockResolutionCache.EXPECT().GetText(chainId, node, key).Return(nil)
			suite.mockAgentEnsRegistryReader.EXPECT().Text(gomock.Any(), chainId, node, key).Return("", nil)
			suite.mockUniversalResolverClient.EXPECT().
				ResolveText(gomock.Any(), chainId, dnsName, innerData).
				Return("", errors.New("OffchainLookup"))
			suite.mockResolutionCache.EXPECT().SetText(chainId, node, key, "")
			suite.mockResolverServiceCodec.EXPECT().EncodeTextResult("").Return([]byte{0x00}, nil)
			suite.mockResponseBuilder.EXPECT().
				Build(gomock.Any(), request.Sender, request.Data, []byte{0x00}).
				Return([]byte{0x42}, nil)

			suite.mockEnsResolutionMetrics.EXPECT().RecordCacheResult(gomock.Any(), "text", false)
			suite.mockEnsResolutionMetrics.EXPECT().RecordFallback(
				gomock.Any(), "text", metrics_interfaces.FallbackOutcomeFailed, chainId)
			suite.mockEnsResolutionMetrics.EXPECT().RecordResolution(
				gomock.Any(), "text", metrics_interfaces.ResolutionSourceAgentEnsRegistry, chainId, gomock.Any(), nil)
			suite.mockEnsResolutionMetrics.EXPECT().RecordSignedResponse(gomock.Any(), nil)
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
			suite.mockResolverServiceCodec.EXPECT().
				Decode(request.Data).
				Return(ens_gateway_interfaces.DecodedQuery{}, errors.New("bad selector"))

			suite.mockEnsResolutionMetrics.EXPECT().RecordDecodeError(gomock.Any(), "undecodable_calldata")
		}

		t.Run("Should return a validation error without touching the cache, registry, or fallback", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingOffchainQueryBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Execute(context.Background(), request)

			var validationErr *apperrors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "data", validationErr.Field)
		})
	})

	t.Run("Given the codec returns an unsupported query kind", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingOffchainQueryTestingSuite) {
			suite.mockResolverServiceCodec.EXPECT().
				Decode(request.Data).
				Return(ens_gateway_interfaces.DecodedQuery{Kind: 99}, nil)

			suite.mockEnsResolutionMetrics.EXPECT().RecordDecodeError(gomock.Any(), "unsupported_query_kind")
		}

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingOffchainQueryBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Execute(context.Background(), request)

			var validationErr *apperrors.ValidationError
			assert.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "data", validationErr.Field)
		})
	})

	t.Run("Given the registry read fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingOffchainQueryTestingSuite) {
			suite.mockResolverServiceCodec.EXPECT().Decode(request.Data).Return(decoded, nil)
			suite.mockResolutionCache.EXPECT().GetText(chainId, node, key).Return(nil)
			suite.mockAgentEnsRegistryReader.EXPECT().
				Text(gomock.Any(), chainId, node, key).
				Return("", errors.New("rpc down"))

			suite.mockEnsResolutionMetrics.EXPECT().RecordCacheResult(gomock.Any(), "text", false)
			suite.mockEnsResolutionMetrics.EXPECT().RecordResolution(
				gomock.Any(),
				"text", metrics_interfaces.ResolutionSourceAgentEnsRegistry, chainId, gomock.Any(), gomock.Any())
		}

		t.Run(
			"Should return the error without caching, invoking the fallback, or building a response",
			func(t *testing.T) {
				t.Parallel()

				suite := WhenResolvingOffchainQueryBeforeEach(t)
				initSuite(suite)

				_, err := suite.sut.Execute(context.Background(), request)

				assert.Error(t, err)
			},
		)
	})

	t.Run("Given the codec fails to encode the result", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingOffchainQueryTestingSuite) {
			suite.mockResolverServiceCodec.EXPECT().Decode(request.Data).Return(decoded, nil)
			suite.mockResolutionCache.EXPECT().GetText(chainId, node, key).Return(nil)
			suite.mockAgentEnsRegistryReader.EXPECT().
				Text(gomock.Any(), chainId, node, key).
				Return("https://agent.example/", nil)
			suite.mockResolutionCache.EXPECT().SetText(chainId, node, key, "https://agent.example/")
			suite.mockResolverServiceCodec.EXPECT().
				EncodeTextResult("https://agent.example/").
				Return(nil, errors.New("encode failed"))

			suite.mockEnsResolutionMetrics.EXPECT().RecordCacheResult(gomock.Any(), "text", false)
			suite.mockEnsResolutionMetrics.EXPECT().RecordResolution(
				gomock.Any(),
				"text", metrics_interfaces.ResolutionSourceAgentEnsRegistry, chainId, gomock.Any(), gomock.Any())
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
			suite.mockResolverServiceCodec.EXPECT().Decode(request.Data).Return(decoded, nil)
			suite.mockResolutionCache.EXPECT().GetText(chainId, node, key).Return(nil)
			suite.mockAgentEnsRegistryReader.EXPECT().
				Text(gomock.Any(), chainId, node, key).
				Return("https://agent.example/", nil)
			suite.mockResolutionCache.EXPECT().SetText(chainId, node, key, "https://agent.example/")
			suite.mockResolverServiceCodec.EXPECT().
				EncodeTextResult("https://agent.example/").
				Return([]byte{0xde, 0xad}, nil)
			suite.mockResponseBuilder.EXPECT().
				Build(gomock.Any(), request.Sender, request.Data, []byte{0xde, 0xad}).
				Return(nil, errors.New("sign failed"))

			suite.mockEnsResolutionMetrics.EXPECT().RecordCacheResult(gomock.Any(), "text", false)
			suite.mockEnsResolutionMetrics.EXPECT().RecordResolution(
				gomock.Any(), "text", metrics_interfaces.ResolutionSourceAgentEnsRegistry, chainId, gomock.Any(), nil)
			suite.mockEnsResolutionMetrics.EXPECT().RecordSignedResponse(gomock.Any(), gomock.Any())
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
			suite.mockResolverServiceCodec.EXPECT().Decode(request.Data).Return(decodedAddr, nil)
			suite.mockResolutionCache.EXPECT().GetAddr(chainId, node).Return(&l2Addr)
			suite.mockResolverServiceCodec.EXPECT().EncodeAddrResult(l2Addr).Return([]byte{0xa1, 0xa2}, nil)
			suite.mockResponseBuilder.EXPECT().
				Build(gomock.Any(), request.Sender, request.Data, []byte{0xa1, 0xa2}).
				Return([]byte{0xb1, 0xb2}, nil)

			suite.mockEnsResolutionMetrics.EXPECT().RecordCacheResult(gomock.Any(), "addr", true)
			suite.mockEnsResolutionMetrics.EXPECT().RecordResolution(
				gomock.Any(), "addr", metrics_interfaces.ResolutionSourceCache, chainId, gomock.Any(), nil)
			suite.mockEnsResolutionMetrics.EXPECT().RecordSignedResponse(gomock.Any(), nil)
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
			suite.mockResolverServiceCodec.EXPECT().Decode(request.Data).Return(decodedAddr, nil)
			suite.mockResolutionCache.EXPECT().GetAddr(chainId, node).Return(nil)
			suite.mockAgentEnsRegistryReader.EXPECT().Addr(gomock.Any(), chainId, node).Return(l2Addr, nil)
			suite.mockResolutionCache.EXPECT().SetAddr(chainId, node, l2Addr)
			suite.mockResolverServiceCodec.EXPECT().EncodeAddrResult(l2Addr).Return([]byte{0xa1, 0xa2}, nil)
			suite.mockResponseBuilder.EXPECT().
				Build(gomock.Any(), request.Sender, request.Data, []byte{0xa1, 0xa2}).
				Return([]byte{0xb1, 0xb2}, nil)

			suite.mockEnsResolutionMetrics.EXPECT().RecordCacheResult(gomock.Any(), "addr", false)
			suite.mockEnsResolutionMetrics.EXPECT().RecordResolution(
				gomock.Any(), "addr", metrics_interfaces.ResolutionSourceAgentEnsRegistry, chainId, gomock.Any(), nil)
			suite.mockEnsResolutionMetrics.EXPECT().RecordSignedResponse(gomock.Any(), nil)
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

	t.Run(
		"Given an addr query and the L2 registry returns the zero address but the L1 fallback resolves",
		func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenResolvingOffchainQueryTestingSuite) {
				suite.mockResolverServiceCodec.EXPECT().Decode(request.Data).Return(decodedAddr, nil)
				suite.mockResolutionCache.EXPECT().GetAddr(chainId, node).Return(nil)
				suite.mockAgentEnsRegistryReader.EXPECT().
					Addr(gomock.Any(), chainId, node).
					Return(common.Address{}, nil)
				suite.mockUniversalResolverClient.EXPECT().
					ResolveAddr(gomock.Any(), chainId, dnsName, addrInnerData).
					Return(l1Addr, nil)
				suite.mockResolutionCache.EXPECT().SetAddr(chainId, node, l1Addr)
				suite.mockResolverServiceCodec.EXPECT().EncodeAddrResult(l1Addr).Return([]byte{0xc1, 0xc2}, nil)
				suite.mockResponseBuilder.EXPECT().
					Build(gomock.Any(), request.Sender, request.Data, []byte{0xc1, 0xc2}).
					Return([]byte{0xd1, 0xd2}, nil)

				suite.mockEnsResolutionMetrics.EXPECT().RecordCacheResult(gomock.Any(), "addr", false)
				suite.mockEnsResolutionMetrics.EXPECT().RecordFallback(
					gomock.Any(), "addr", metrics_interfaces.FallbackOutcomeResolved, chainId)
				suite.mockEnsResolutionMetrics.EXPECT().RecordResolution(
					gomock.Any(),
					"addr", metrics_interfaces.ResolutionSourceUniversalResolver, chainId, gomock.Any(), nil)
				suite.mockEnsResolutionMetrics.EXPECT().RecordSignedResponse(gomock.Any(), nil)
			}

			t.Run("Should sign the L1 fallback address", func(t *testing.T) {
				t.Parallel()

				suite := WhenResolvingOffchainQueryBeforeEach(t)
				initSuite(suite)

				response, err := suite.sut.Execute(context.Background(), request)

				assert.NoError(t, err)
				assert.Equal(t, "0xd1d2", response.Data)
			})
		},
	)

	t.Run("Given an addr query and both L2 and L1 return the zero address", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingOffchainQueryTestingSuite) {
			suite.mockResolverServiceCodec.EXPECT().Decode(request.Data).Return(decodedAddr, nil)
			suite.mockResolutionCache.EXPECT().GetAddr(chainId, node).Return(nil)
			suite.mockAgentEnsRegistryReader.EXPECT().Addr(gomock.Any(), chainId, node).Return(common.Address{}, nil)
			suite.mockUniversalResolverClient.EXPECT().
				ResolveAddr(gomock.Any(), chainId, dnsName, addrInnerData).
				Return(common.Address{}, nil)
			suite.mockResolutionCache.EXPECT().SetAddr(chainId, node, common.Address{})
			suite.mockResolverServiceCodec.EXPECT().EncodeAddrResult(common.Address{}).Return([]byte{0x00}, nil)
			suite.mockResponseBuilder.EXPECT().
				Build(gomock.Any(), request.Sender, request.Data, []byte{0x00}).
				Return([]byte{0x42}, nil)

			suite.mockEnsResolutionMetrics.EXPECT().RecordCacheResult(gomock.Any(), "addr", false)
			suite.mockEnsResolutionMetrics.EXPECT().RecordFallback(
				gomock.Any(), "addr", metrics_interfaces.FallbackOutcomeEmpty, chainId)
			suite.mockEnsResolutionMetrics.EXPECT().RecordResolution(
				gomock.Any(), "addr", metrics_interfaces.ResolutionSourceAgentEnsRegistry, chainId, gomock.Any(), nil)
			suite.mockEnsResolutionMetrics.EXPECT().RecordSignedResponse(gomock.Any(), nil)
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
			suite.mockResolverServiceCodec.EXPECT().Decode(request.Data).Return(decodedAddr, nil)
			suite.mockResolutionCache.EXPECT().GetAddr(chainId, node).Return(nil)
			suite.mockAgentEnsRegistryReader.EXPECT().Addr(gomock.Any(), chainId, node).Return(common.Address{}, nil)
			suite.mockUniversalResolverClient.EXPECT().
				ResolveAddr(gomock.Any(), chainId, dnsName, addrInnerData).
				Return(common.Address{}, errors.New("OffchainLookup"))
			suite.mockResolutionCache.EXPECT().SetAddr(chainId, node, common.Address{})
			suite.mockResolverServiceCodec.EXPECT().EncodeAddrResult(common.Address{}).Return([]byte{0x00}, nil)
			suite.mockResponseBuilder.EXPECT().
				Build(gomock.Any(), request.Sender, request.Data, []byte{0x00}).
				Return([]byte{0x42}, nil)

			suite.mockEnsResolutionMetrics.EXPECT().RecordCacheResult(gomock.Any(), "addr", false)
			suite.mockEnsResolutionMetrics.EXPECT().RecordFallback(
				gomock.Any(), "addr", metrics_interfaces.FallbackOutcomeFailed, chainId)
			suite.mockEnsResolutionMetrics.EXPECT().RecordResolution(
				gomock.Any(), "addr", metrics_interfaces.ResolutionSourceAgentEnsRegistry, chainId, gomock.Any(), nil)
			suite.mockEnsResolutionMetrics.EXPECT().RecordSignedResponse(gomock.Any(), nil)
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
			suite.mockResolverServiceCodec.EXPECT().Decode(request.Data).Return(decodedAddr, nil)
			suite.mockResolutionCache.EXPECT().GetAddr(chainId, node).Return(nil)
			suite.mockAgentEnsRegistryReader.EXPECT().
				Addr(gomock.Any(), chainId, node).
				Return(common.Address{}, errors.New("rpc down"))

			suite.mockEnsResolutionMetrics.EXPECT().RecordCacheResult(gomock.Any(), "addr", false)
			suite.mockEnsResolutionMetrics.EXPECT().RecordResolution(
				gomock.Any(),
				"addr", metrics_interfaces.ResolutionSourceAgentEnsRegistry, chainId, gomock.Any(), gomock.Any())
		}

		t.Run(
			"Should return the error without caching, invoking the fallback, or building a response",
			func(t *testing.T) {
				t.Parallel()

				suite := WhenResolvingOffchainQueryBeforeEach(t)
				initSuite(suite)

				_, err := suite.sut.Execute(context.Background(), request)

				assert.Error(t, err)
			},
		)
	})
}

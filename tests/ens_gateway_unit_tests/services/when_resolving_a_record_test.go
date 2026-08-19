package services_test

import (
	"context"
	"testing"

	gateway_services "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/services"
	ens_gateway_interfaces "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/services/interfaces"
	traces_interfaces "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/traces/interfaces"
	"github.com/6022-labs/agentic-ens-gateway/tests/common_tests/model_builders"
	metrics_mocks "github.com/6022-labs/agentic-ens-gateway/tests/ens_gateway_mocks/metrics_mocks/interfaces_mocks"
	services_mocks "github.com/6022-labs/agentic-ens-gateway/tests/ens_gateway_mocks/services_mocks/interfaces_mocks"
	traces_mocks "github.com/6022-labs/agentic-ens-gateway/tests/ens_gateway_mocks/traces_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenResolvingARecordTestingSuite struct {
	sut *gateway_services.RecordResolver

	mockResolverServiceCodec    *services_mocks.MockResolverServiceCodecInterface
	mockAgentEnsRegistryReader  *services_mocks.MockAgentEnsRegistryReaderInterface
	mockUniversalResolverClient *services_mocks.MockUniversalResolverClientInterface
	mockResolutionCache         *services_mocks.MockResolutionCacheInterface
	mockEnsResolutionMetrics    *metrics_mocks.MockEnsResolutionMetricsInterface
	mockSpan                    *traces_mocks.MockSpan
}

func WhenResolvingARecordBeforeEach(t *testing.T) *WhenResolvingARecordTestingSuite {
	mockController := gomock.NewController(t)

	mockResolverServiceCodec := services_mocks.NewMockResolverServiceCodecInterface(mockController)
	mockAgentEnsRegistryReader := services_mocks.NewMockAgentEnsRegistryReaderInterface(mockController)
	mockUniversalResolverClient := services_mocks.NewMockUniversalResolverClientInterface(mockController)
	mockResolutionCache := services_mocks.NewMockResolutionCacheInterface(mockController)
	mockEnsResolutionMetrics := metrics_mocks.NewMockEnsResolutionMetricsInterface(mockController)
	mockResolutionTracer := traces_mocks.NewMockResolutionTracerInterface(mockController)

	mockSpan := traces_mocks.NewMockSpan(mockController)
	mockSpan.EXPECT().End().AnyTimes()
	mockSpan.EXPECT().Fail(gomock.Any()).AnyTimes()
	mockResolutionTracer.EXPECT().StartRegistryRead(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, _ uint64) (context.Context, traces_interfaces.Span) {
			return ctx, mockSpan
		}).AnyTimes()
	mockResolutionTracer.EXPECT().StartUniversalResolverFallback(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, _ uint64) (context.Context, traces_interfaces.Span) {
			return ctx, mockSpan
		}).AnyTimes()

	sut := gateway_services.NewRecordResolver(
		zap.NewNop(),
		mockResolverServiceCodec,
		mockAgentEnsRegistryReader,
		mockUniversalResolverClient,
		mockResolutionCache,
		mockEnsResolutionMetrics,
		mockResolutionTracer,
	)

	return &WhenResolvingARecordTestingSuite{
		sut:                         sut,
		mockResolverServiceCodec:    mockResolverServiceCodec,
		mockAgentEnsRegistryReader:  mockAgentEnsRegistryReader,
		mockUniversalResolverClient: mockUniversalResolverClient,
		mockResolutionCache:         mockResolutionCache,
		mockEnsResolutionMetrics:    mockEnsResolutionMetrics,
		mockSpan:                    mockSpan,
	}
}

func TestWhenResolvingARecord(t *testing.T) {
	t.Parallel()

	t.Run("Given the text record is cached", func(t *testing.T) {
		t.Parallel()

		query := model_builders.NewDecodedQueryBuilder().Build()
		cached := "https://agent.example/"

		initSuite := func(suite *WhenResolvingARecordTestingSuite) {
			suite.mockResolutionCache.EXPECT().GetText(query.ChainId, query.Node, query.Key).Return(&cached)
			suite.mockEnsResolutionMetrics.EXPECT().RecordCacheResult(gomock.Any(), "text", true)
			suite.mockResolverServiceCodec.EXPECT().EncodeTextResult(cached).Return([]byte{0xde, 0xad}, nil)
		}

		t.Run("Should answer from the cache without reading the registry", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingARecordBeforeEach(t)
			initSuite(suite)

			result, source, err := suite.sut.Resolve(context.Background(), query)

			assert.NoError(t, err)
			assert.Equal(t, []byte{0xde, 0xad}, result)
			assert.Equal(t, "cache", source)
		})
	})

	t.Run("Given the registry answers empty and the L1 fallback resolves", func(t *testing.T) {
		t.Parallel()

		query := model_builders.NewDecodedQueryBuilder().Build()

		initSuite := func(suite *WhenResolvingARecordTestingSuite) {
			suite.mockResolutionCache.EXPECT().GetText(query.ChainId, query.Node, query.Key).Return(nil)
			suite.mockEnsResolutionMetrics.EXPECT().RecordCacheResult(gomock.Any(), "text", false)
			suite.mockAgentEnsRegistryReader.EXPECT().
				Text(gomock.Any(), query.ChainId, query.Node, query.Key).
				Return("", nil)
			suite.mockUniversalResolverClient.EXPECT().
				ResolveText(gomock.Any(), query.ChainId, query.DnsName, query.InnerData).
				Return("from-l1", nil)
			suite.mockEnsResolutionMetrics.EXPECT().
				RecordFallback(gomock.Any(), "text", "resolved", query.ChainId)
			suite.mockResolutionCache.EXPECT().SetText(query.ChainId, query.Node, query.Key, "from-l1")
			suite.mockResolverServiceCodec.EXPECT().EncodeTextResult("from-l1").Return([]byte{0x01}, nil)
		}

		t.Run("Should return the fallback value and report the universal resolver as the source", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingARecordBeforeEach(t)
			initSuite(suite)

			result, source, err := suite.sut.Resolve(context.Background(), query)

			assert.NoError(t, err)
			assert.Equal(t, []byte{0x01}, result)
			assert.Equal(t, "universal_resolver", source)
		})
	})

	t.Run("Given the registry read fails", func(t *testing.T) {
		t.Parallel()

		query := model_builders.NewDecodedQueryBuilder().
			WithKind(ens_gateway_interfaces.QueryKindAddr).
			Build()

		initSuite := func(suite *WhenResolvingARecordTestingSuite) {
			suite.mockResolutionCache.EXPECT().GetAddr(query.ChainId, query.Node).Return(nil)
			suite.mockEnsResolutionMetrics.EXPECT().RecordCacheResult(gomock.Any(), "addr", false)
			suite.mockAgentEnsRegistryReader.EXPECT().
				Addr(gomock.Any(), query.ChainId, query.Node).
				Return(common.Address{}, assert.AnError)
		}

		t.Run("Should surface the error without consulting the L1 fallback", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingARecordBeforeEach(t)
			initSuite(suite)

			result, _, err := suite.sut.Resolve(context.Background(), query)

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})
}

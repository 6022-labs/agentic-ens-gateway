package ens_gateway_mvc_unit_tests_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	ens_gateway_interfaces "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/services/interfaces"
	traces_interfaces "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/traces/interfaces"
	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/use_cases"
	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_mvc"
	metrics_mocks "github.com/6022-labs/agentic-ens-gateway/tests/ens_gateway_mocks/metrics_mocks/interfaces_mocks"
	services_mocks "github.com/6022-labs/agentic-ens-gateway/tests/ens_gateway_mocks/services_mocks/interfaces_mocks"
	traces_mocks "github.com/6022-labs/agentic-ens-gateway/tests/ens_gateway_mocks/traces_mocks/interfaces_mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

const (
	testSender  = "0x1111111111111111111111111111111111111111"
	testDataHex = "0x010203"
)

type WhenHandlingCcipReadRequestTestingSuite struct {
	app *fiber.App

	mockResolverServiceCodec    *services_mocks.MockResolverServiceCodecInterface
	mockAgentEnsRegistryReader  *services_mocks.MockAgentEnsRegistryReaderInterface
	mockUniversalResolverClient *services_mocks.MockUniversalResolverClientInterface
	mockResponseBuilder         *services_mocks.MockResponseBuilderInterface
	mockResolutionCache         *services_mocks.MockResolutionCacheInterface
	mockEnsResolutionMetrics    *metrics_mocks.MockEnsResolutionMetricsInterface
	mockResolutionTracer        *traces_mocks.MockResolutionTracerInterface
	mockResolutionSpan          *traces_mocks.MockResolutionSpan
}

func WhenHandlingCcipReadRequestBeforeEach(t *testing.T) *WhenHandlingCcipReadRequestTestingSuite {
	mockController := gomock.NewController(t)

	mockResolverServiceCodec := services_mocks.NewMockResolverServiceCodecInterface(mockController)
	mockAgentEnsRegistryReader := services_mocks.NewMockAgentEnsRegistryReaderInterface(mockController)
	mockUniversalResolverClient := services_mocks.NewMockUniversalResolverClientInterface(mockController)
	mockResponseBuilder := services_mocks.NewMockResponseBuilderInterface(mockController)
	mockResolutionCache := services_mocks.NewMockResolutionCacheInterface(mockController)
	mockEnsResolutionMetrics := metrics_mocks.NewMockEnsResolutionMetricsInterface(mockController)
	mockResolutionTracer := traces_mocks.NewMockResolutionTracerInterface(mockController)
	mockResolutionSpan := traces_mocks.NewMockResolutionSpan(mockController)

	resolveOffchainQuery := use_cases.NewResolveOffchainQuery(
		zap.NewNop(),
		mockResolverServiceCodec,
		mockAgentEnsRegistryReader,
		mockUniversalResolverClient,
		mockResponseBuilder,
		mockResolutionCache,
		mockEnsResolutionMetrics,
		mockResolutionTracer,
	)

	controller := ens_gateway_mvc.NewCcipReadController(zap.NewNop(), resolveOffchainQuery)
	app := fiber.New()
	controller.RegisterRoutes(app)

	return &WhenHandlingCcipReadRequestTestingSuite{
		app:                         app,
		mockResolverServiceCodec:    mockResolverServiceCodec,
		mockAgentEnsRegistryReader:  mockAgentEnsRegistryReader,
		mockUniversalResolverClient: mockUniversalResolverClient,
		mockResponseBuilder:         mockResponseBuilder,
		mockResolutionCache:         mockResolutionCache,
		mockEnsResolutionMetrics:    mockEnsResolutionMetrics,
		mockResolutionTracer:        mockResolutionTracer,
		mockResolutionSpan:          mockResolutionSpan,
	}
}

// expectResolvedText drives the shortest path: a cached text answer handed to the builder,
// which returns signed or fails with buildErr.
func expectResolvedText(suite *WhenHandlingCcipReadRequestTestingSuite, signed []byte, buildErr error) {
	node := [32]byte{0x01}
	cached := "https://agent.example/"

	suite.mockResolverServiceCodec.EXPECT().
		Decode(gomock.Any()).
		Return(ens_gateway_interfaces.DecodedQuery{
			Kind: ens_gateway_interfaces.QueryKindText, ChainId: 80002, Node: node, Key: "url",
		}, nil)

	suite.mockResolutionTracer.EXPECT().
		StartResolution(gomock.Any(), "text", uint64(80002)).
		DoAndReturn(func(ctx context.Context, _ string, _ uint64) (context.Context, traces_interfaces.ResolutionSpan) {
			return ctx, suite.mockResolutionSpan
		})
	suite.mockResolutionSpan.EXPECT().RecordSource(gomock.Any()).AnyTimes()
	suite.mockResolutionSpan.EXPECT().Fail(gomock.Any()).AnyTimes()
	suite.mockResolutionSpan.EXPECT().End().AnyTimes()

	suite.mockResolutionCache.EXPECT().GetText(uint64(80002), node, "url").Return(&cached)
	suite.mockEnsResolutionMetrics.EXPECT().RecordCacheResult(gomock.Any(), "text", true)
	suite.mockResolverServiceCodec.EXPECT().EncodeTextResult(cached).Return([]byte{0xde, 0xad}, nil)
	suite.mockEnsResolutionMetrics.EXPECT().
		RecordResolution(gomock.Any(), "text", gomock.Any(), uint64(80002), gomock.Any(), gomock.Any())

	suite.mockResponseBuilder.EXPECT().
		Build(gomock.Any(), gomock.Any(), gomock.Any(), []byte{0xde, 0xad}).
		Return(signed, buildErr)
	suite.mockEnsResolutionMetrics.EXPECT().RecordSignedResponse(gomock.Any(), buildErr)
}

func TestWhenHandlingCcipReadRequest(t *testing.T) {
	t.Parallel()

	t.Run("Given a valid GET request with a .json suffix", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingCcipReadRequestTestingSuite) {
			expectResolvedText(suite, []byte{0xbe, 0xef}, nil)
		}

		t.Run("Should strip .json, call the use case, and return 200 with the data", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingCcipReadRequestBeforeEach(t)
			initSuite(suite)

			req := httptest.NewRequest(http.MethodGet, "/"+testSender+"/"+testDataHex+".json", nil)
			resp, err := suite.app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			body, _ := io.ReadAll(resp.Body)
			assert.JSONEq(t, `{"data":"0xbeef"}`, string(body))
		})
	})

	t.Run("Given a valid POST request", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingCcipReadRequestTestingSuite) {
			expectResolvedText(suite, []byte{0xbe, 0xef}, nil)
		}

		t.Run("Should call the use case and return 200", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingCcipReadRequestBeforeEach(t)
			initSuite(suite)

			req := httptest.NewRequest(
				http.MethodPost,
				"/"+testSender,
				strings.NewReader(`{"data":"`+testDataHex+`"}`),
			)
			req.Header.Set("Content-Type", "application/json")
			resp, err := suite.app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	})

	t.Run("Given malformed call data hex", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return 400 without reaching the codec", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingCcipReadRequestBeforeEach(t)

			req := httptest.NewRequest(http.MethodGet, "/"+testSender+"/not-hex.json", nil)
			resp, err := suite.app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})

	t.Run("Given an invalid sender address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return 400 without reaching the codec", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingCcipReadRequestBeforeEach(t)

			req := httptest.NewRequest(http.MethodGet, "/not-an-address/"+testDataHex+".json", nil)
			resp, err := suite.app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})

	t.Run("Given a POST request with an invalid body", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return 400 without reaching the codec", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingCcipReadRequestBeforeEach(t)

			req := httptest.NewRequest(
				http.MethodPost,
				"/"+testSender,
				strings.NewReader("not json"),
			)
			req.Header.Set("Content-Type", "application/json")
			resp, err := suite.app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})

	t.Run("Given the use case fails with a validation error", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingCcipReadRequestTestingSuite) {
			suite.mockResolverServiceCodec.EXPECT().
				Decode(gomock.Any()).
				Return(ens_gateway_interfaces.DecodedQuery{}, errors.New("invalid resolver call data"))
			suite.mockEnsResolutionMetrics.EXPECT().RecordDecodeError(gomock.Any(), "undecodable_calldata")
		}

		t.Run("Should return 400 with the field and message", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingCcipReadRequestBeforeEach(t)
			initSuite(suite)

			req := httptest.NewRequest(http.MethodGet, "/"+testSender+"/"+testDataHex+".json", nil)
			resp, err := suite.app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

			var body map[string]string
			assert.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
			assert.Equal(t, "data", body["field"])
			assert.Equal(t, "invalid resolver call data: invalid resolver call data", body["message"])
		})
	})

	t.Run("Given the use case fails with an untyped error", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingCcipReadRequestTestingSuite) {
			expectResolvedText(suite, nil, errors.New("boom"))
		}

		t.Run("Should return 500", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingCcipReadRequestBeforeEach(t)
			initSuite(suite)

			req := httptest.NewRequest(http.MethodGet, "/"+testSender+"/"+testDataHex+".json", nil)
			resp, err := suite.app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		})
	})
}

package services_test

import (
	"context"
	"testing"

	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_blockchain/factory"
	blockchain_services "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_blockchain/services"
	blockchain_settings "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_blockchain/settings"
	blockchain_mocks "github.com/6022-labs/agentic-ens-gateway/tests/ens_gateway_blockchain_mocks/metrics_mocks/interfaces_mocks"
	services_mocks "github.com/6022-labs/agentic-ens-gateway/tests/ens_gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenResolvingTextThroughTheUniversalResolverTestingSuite struct {
	sut *blockchain_services.UniversalResolverClient

	mockResolverServiceCodec *services_mocks.MockResolverServiceCodecInterface
}

func WhenResolvingTextThroughTheUniversalResolverBeforeEach(
	t *testing.T,
	universalResolverSettings *blockchain_settings.UniversalResolverSettings,
) *WhenResolvingTextThroughTheUniversalResolverTestingSuite {
	mockController := gomock.NewController(t)
	logger := zap.NewNop()

	chainsSettings := blockchain_settings.ChainsSettings{}
	mockResolverServiceCodec := services_mocks.NewMockResolverServiceCodecInterface(mockController)

	sut := blockchain_services.NewUniversalResolverClient(
		logger,
		&chainsSettings,
		universalResolverSettings,
		mockResolverServiceCodec,
		factory.NewEthClientFactory(
			&chainsSettings,
			blockchain_mocks.NewMockBlockchainRpcMetricsInterface(mockController),
		),
	)

	return &WhenResolvingTextThroughTheUniversalResolverTestingSuite{
		sut:                      sut,
		mockResolverServiceCodec: mockResolverServiceCodec,
	}
}

func TestWhenResolvingTextThroughTheUniversalResolver(t *testing.T) {
	t.Parallel()

	t.Run("Given the source chain has no mapping at all", func(t *testing.T) {
		t.Parallel()

		universalResolverSettings := &blockchain_settings.UniversalResolverSettings{
			Resolvers: map[uint64]common.Address{},
			Mappings:  map[uint64]uint64{},
		}

		t.Run("Should return an empty value without calling the codec", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingTextThroughTheUniversalResolverBeforeEach(t, universalResolverSettings)

			value, err := suite.sut.ResolveText(context.Background(), 80002, []byte{0x01}, []byte{0x02})

			assert.NoError(t, err)
			assert.Equal(t, "", value)
		})
	})

	t.Run("Given the mapped resolver chain has no configured rpc url", func(t *testing.T) {
		t.Parallel()

		universalResolverSettings := &blockchain_settings.UniversalResolverSettings{
			Resolvers: map[uint64]common.Address{1: common.HexToAddress("0x1111111111111111111111111111111111111111")},
			Mappings:  map[uint64]uint64{80002: 1},
		}

		t.Run("Should leave the fallback disabled and return an empty value", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingTextThroughTheUniversalResolverBeforeEach(t, universalResolverSettings)

			value, err := suite.sut.ResolveText(context.Background(), 80002, []byte{0x01}, []byte{0x02})

			assert.NoError(t, err)
			assert.Equal(t, "", value)
		})
	})
}

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

type WhenResolvingAddrThroughTheUniversalResolverTestingSuite struct {
	sut *blockchain_services.UniversalResolverClient

	mockResolverServiceCodec *services_mocks.MockResolverServiceCodecInterface
}

func WhenResolvingAddrThroughTheUniversalResolverBeforeEach(
	t *testing.T,
	universalResolverSettings *blockchain_settings.UniversalResolverSettings,
) *WhenResolvingAddrThroughTheUniversalResolverTestingSuite {
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

	return &WhenResolvingAddrThroughTheUniversalResolverTestingSuite{
		sut:                      sut,
		mockResolverServiceCodec: mockResolverServiceCodec,
	}
}

func TestWhenResolvingAddrThroughTheUniversalResolver(t *testing.T) {
	t.Parallel()

	t.Run("Given the source chain has no mapping at all", func(t *testing.T) {
		t.Parallel()

		universalResolverSettings := &blockchain_settings.UniversalResolverSettings{
			Resolvers: map[uint64]common.Address{},
			Mappings:  map[uint64]uint64{},
		}

		t.Run("Should return the zero address without calling the codec", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingAddrThroughTheUniversalResolverBeforeEach(t, universalResolverSettings)

			value, err := suite.sut.ResolveAddr(context.Background(), 80002, []byte{0x01}, []byte{0x02})

			assert.NoError(t, err)
			assert.Equal(t, common.Address{}, value)
		})
	})

	t.Run("Given the mapped resolver chain has no configured rpc url", func(t *testing.T) {
		t.Parallel()

		universalResolverSettings := &blockchain_settings.UniversalResolverSettings{
			Resolvers: map[uint64]common.Address{1: common.HexToAddress("0x1111111111111111111111111111111111111111")},
			Mappings:  map[uint64]uint64{80002: 1},
		}

		t.Run("Should leave the fallback disabled and return the zero address", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingAddrThroughTheUniversalResolverBeforeEach(t, universalResolverSettings)

			value, err := suite.sut.ResolveAddr(context.Background(), 80002, []byte{0x01}, []byte{0x02})

			assert.NoError(t, err)
			assert.Equal(t, common.Address{}, value)
		})
	})
}

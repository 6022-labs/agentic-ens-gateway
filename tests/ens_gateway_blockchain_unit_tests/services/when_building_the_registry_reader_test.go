package services_test

import (
	"context"
	"testing"

	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_blockchain/factory"
	blockchain_services "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_blockchain/services"
	blockchain_settings "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_blockchain/settings"
	blockchain_mocks "github.com/6022-labs/agentic-ens-gateway/tests/ens_gateway_blockchain_mocks/metrics_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenBuildingTheRegistryReaderTestingSuite struct {
	sut *blockchain_services.AgentEnsRegistryReader
}

func WhenBuildingTheRegistryReaderBeforeEach(t *testing.T) *WhenBuildingTheRegistryReaderTestingSuite {
	logger := zap.NewNop()

	chainsSettings := blockchain_settings.ChainsSettings{
		80002: &blockchain_settings.ChainConfig{
			RpcHttpUrl:            "https://rpc.invalid/amoy",
			RequestTimeoutSeconds: 10,
		},
	}

	registries := blockchain_settings.AgentEnsRegistriesSettings{
		137:   common.HexToAddress("0x23dbc335f966d9869940FB55f61E919118AAD236"),
		80002: common.HexToAddress("0x37c5879529E5b0241663ed40feeebb2326BBEE0d"),
	}

	sut := blockchain_services.NewAgentEnsRegistryReader(
		logger,
		&chainsSettings,
		&registries,
		factory.NewEthClientFactory(
			&chainsSettings,
			blockchain_mocks.NewMockBlockchainRpcMetricsInterface(gomock.NewController(t)),
		),
	)

	return &WhenBuildingTheRegistryReaderTestingSuite{sut: sut}
}

func TestWhenBuildingTheRegistryReader(t *testing.T) {
	t.Parallel()

	t.Run("Given a baked-in registry chain has no configured rpc url", func(t *testing.T) {
		t.Parallel()

		t.Run("Should skip that chain instead of aborting the process", func(t *testing.T) {
			t.Parallel()

			suite := WhenBuildingTheRegistryReaderBeforeEach(t)

			value, err := suite.sut.Text(context.Background(), 137, [32]byte{0x01}, "url")

			assert.NoError(t, err)
			assert.Equal(t, "", value)
		})
	})
}

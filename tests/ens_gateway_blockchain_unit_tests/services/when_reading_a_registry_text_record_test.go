package services_test

import (
	"context"
	"testing"

	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_blockchain/factory"
	blockchain_services "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_blockchain/services"
	blockchain_settings "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_blockchain/settings"
	blockchain_mocks "github.com/6022-labs/agentic-ens-gateway/tests/ens_gateway_blockchain_mocks/metrics_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenReadingARegistryTextRecordTestingSuite struct {
	sut *blockchain_services.AgentEnsRegistryReader
}

func WhenReadingARegistryTextRecordBeforeEach(t *testing.T) *WhenReadingARegistryTextRecordTestingSuite {
	logger := zap.NewNop()
	chainsSettings := blockchain_settings.ChainsSettings{}
	registries := blockchain_settings.AgentEnsRegistriesSettings{}

	sut := blockchain_services.NewAgentEnsRegistryReader(
		logger,
		&chainsSettings,
		&registries,
		factory.NewEthClientFactory(
			&chainsSettings,
			blockchain_mocks.NewMockBlockchainRpcMetricsInterface(gomock.NewController(t)),
		),
	)

	return &WhenReadingARegistryTextRecordTestingSuite{sut: sut}
}

func TestWhenReadingARegistryTextRecord(t *testing.T) {
	t.Parallel()

	t.Run("Given the chain has no bound AgentEnsRegistry", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an empty value without an error so the caller falls back to L1", func(t *testing.T) {
			t.Parallel()

			suite := WhenReadingARegistryTextRecordBeforeEach(t)

			value, err := suite.sut.Text(context.Background(), 80002, [32]byte{0x01}, "url")

			assert.NoError(t, err)
			assert.Equal(t, "", value)
		})
	})
}

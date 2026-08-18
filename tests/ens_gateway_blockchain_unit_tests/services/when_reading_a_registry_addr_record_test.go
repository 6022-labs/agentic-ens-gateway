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

type WhenReadingARegistryAddrRecordTestingSuite struct {
	sut *blockchain_services.AgentEnsRegistryReader
}

func WhenReadingARegistryAddrRecordBeforeEach(t *testing.T) *WhenReadingARegistryAddrRecordTestingSuite {
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

	return &WhenReadingARegistryAddrRecordTestingSuite{sut: sut}
}

func TestWhenReadingARegistryAddrRecord(t *testing.T) {
	t.Parallel()

	t.Run("Given the chain has no bound AgentEnsRegistry", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the zero address without an error so the caller falls back to L1", func(t *testing.T) {
			t.Parallel()

			suite := WhenReadingARegistryAddrRecordBeforeEach(t)

			value, err := suite.sut.Addr(context.Background(), 80002, [32]byte{0x01})

			assert.NoError(t, err)
			assert.Equal(t, common.Address{}, value)
		})
	})
}

package factory_test

import (
	"testing"

	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_blockchain/factory"
	blockchain_settings "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_blockchain/settings"
	blockchain_mocks "github.com/6022-labs/agentic-ens-gateway/tests/ens_gateway_blockchain_mocks/metrics_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenGettingAnEthClientTestingSuite struct {
	sut *factory.EthClientFactory
}

func WhenGettingAnEthClientBeforeEach(
	t *testing.T,
	chainsSettings blockchain_settings.ChainsSettings,
) *WhenGettingAnEthClientTestingSuite {
	mockController := gomock.NewController(t)

	sut := factory.NewEthClientFactory(
		&chainsSettings,
		blockchain_mocks.NewMockBlockchainRpcMetricsInterface(mockController),
	)

	return &WhenGettingAnEthClientTestingSuite{sut: sut}
}

func TestWhenGettingAnEthClient(t *testing.T) {
	t.Parallel()

	t.Run("Given the chain is not configured", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error rather than a nil client", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAnEthClientBeforeEach(t, blockchain_settings.ChainsSettings{})

			client, err := suite.sut.For(80002)

			assert.Error(t, err)
			assert.Nil(t, client)
		})
	})
}

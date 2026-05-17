package settings_test

import (
	"testing"

	blockchain_settings "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_blockchain/settings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingUniversalResolverSettings(t *testing.T) {
	t.Parallel()

	canonicalResolver := common.HexToAddress("0xeEeEEEeE14D718C2B47D9923Deab1335E144EeEe").Hex()

	t.Run("Given the embedded universal_resolvers.json", func(t *testing.T) {
		t.Parallel()

		s := blockchain_settings.NewUniversalResolverSettings(zap.NewNop())

		t.Run("Should expose UniversalResolver addresses for both mainnet and sepolia", func(t *testing.T) {
			t.Parallel()
			mainnet, ok := s.Resolvers[1]
			assert.True(t, ok)
			assert.Equal(t, canonicalResolver, mainnet.Hex())

			sepolia, ok := s.Resolvers[11155111]
			assert.True(t, ok)
			assert.Equal(t, canonicalResolver, sepolia.Hex())
		})

		t.Run("Should map known source chain ids to their UniversalResolver chain ids", func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, uint64(1), s.Mappings[1])
			assert.Equal(t, uint64(1), s.Mappings[137])
			assert.Equal(t, uint64(11155111), s.Mappings[80002])
		})

		t.Run("Get should compose the two maps for a configured source chain", func(t *testing.T) {
			t.Parallel()
			resolverChainId, addr, ok := s.Get(137)
			assert.True(t, ok)
			assert.Equal(t, uint64(1), resolverChainId)
			assert.Equal(t, canonicalResolver, addr.Hex())
		})

		t.Run("Get should resolve chain 1 to itself", func(t *testing.T) {
			t.Parallel()
			resolverChainId, addr, ok := s.Get(1)
			assert.True(t, ok)
			assert.Equal(t, uint64(1), resolverChainId)
			assert.Equal(t, canonicalResolver, addr.Hex())
		})

		t.Run("Get should return ok=false for an unmapped source chain id", func(t *testing.T) {
			t.Parallel()
			_, _, ok := s.Get(999)
			assert.False(t, ok)
		})
	})

	t.Run("Given a source chain is mapped to a resolver chain that has no configured address", func(t *testing.T) {
		t.Parallel()

		t.Run("Get should return the resolver chain id but ok=false", func(t *testing.T) {
			t.Parallel()

			s := &blockchain_settings.UniversalResolverSettings{
				Resolvers: map[uint64]common.Address{},
				Mappings:  map[uint64]uint64{42: 999},
			}

			resolverChainId, _, ok := s.Get(42)
			assert.False(t, ok)
			assert.Equal(t, uint64(999), resolverChainId)
		})
	})
}

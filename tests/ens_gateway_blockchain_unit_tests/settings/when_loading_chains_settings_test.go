package settings_test

import (
	"testing"

	blockchain_settings "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_blockchain/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestWhenLoadingChainsSettings(t *testing.T) {
	t.Parallel()

	t.Run("Given chains are provided as a map keyed by chain id", func(t *testing.T) {
		t.Parallel()

		t.Run("Should populate every entry", func(t *testing.T) {
			t.Parallel()

			k := koanf.New(".")
			require.NoError(t, k.Load(confmap.Provider(map[string]any{
				"chains.137.rpc_http_url":   "https://polygon-rpc.com",
				"chains.80002.rpc_http_url": "https://polygon-amoy.example",
			}, "."), nil))

			s := blockchain_settings.NewChainsSettings(zap.NewNop(), k)
			require.NotNil(t, s)
			loaded := *s
			require.Len(t, loaded, 2)

			require.NotNil(t, loaded[137])
			assert.Equal(t, "https://polygon-rpc.com", loaded[137].RpcHttpUrl)

			require.NotNil(t, loaded[80002])
			assert.Equal(t, "https://polygon-amoy.example", loaded[80002].RpcHttpUrl)
		})
	})

	t.Run("Given the rpc url has surrounding whitespace", func(t *testing.T) {
		t.Parallel()

		t.Run("Should trim it", func(t *testing.T) {
			t.Parallel()

			k := koanf.New(".")
			require.NoError(t, k.Load(confmap.Provider(map[string]any{
				"chains.137.rpc_http_url": "  https://polygon-rpc.com  ",
			}, "."), nil))

			s := blockchain_settings.NewChainsSettings(zap.NewNop(), k)
			loaded := *s

			assert.Equal(t, "https://polygon-rpc.com", loaded[137].RpcHttpUrl)
		})
	})

	t.Run("Given Get is called", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the matching chain or nil for an unknown id", func(t *testing.T) {
			t.Parallel()

			k := koanf.New(".")
			require.NoError(t, k.Load(confmap.Provider(map[string]any{
				"chains.137.rpc_http_url": "https://polygon-rpc.com",
			}, "."), nil))

			s := blockchain_settings.NewChainsSettings(zap.NewNop(), k)

			match := s.Get(137)
			require.NotNil(t, match)
			assert.Equal(t, "https://polygon-rpc.com", match.RpcHttpUrl)
			assert.Nil(t, s.Get(999))
		})
	})

	t.Run("Given a chain without an explicit request timeout", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"chains.137.rpc_http_url": "https://polygon-rpc.com",
		}, "."), nil)

		t.Run("Should default it so no RPC call can hang unbounded", func(t *testing.T) {
			t.Parallel()

			result := blockchain_settings.NewChainsSettings(zap.NewNop(), k)

			assert.Equal(t, uint64(10), result.Get(137).RequestTimeoutSeconds)
		})
	})

	t.Run("Given a chain with an explicit request timeout", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"chains.137.rpc_http_url":            "https://polygon-rpc.com",
			"chains.137.request_timeout_seconds": 3,
		}, "."), nil)

		t.Run("Should keep the configured value", func(t *testing.T) {
			t.Parallel()

			result := blockchain_settings.NewChainsSettings(zap.NewNop(), k)

			assert.Equal(t, uint64(3), result.Get(137).RequestTimeoutSeconds)
		})
	})
}

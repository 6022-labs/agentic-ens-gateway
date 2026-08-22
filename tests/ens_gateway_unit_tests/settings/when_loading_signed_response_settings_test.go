package settings_test

import (
	"testing"
	"time"

	gateway_settings "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// Throwaway key, never a deployed signer.
const testPrivateKey = "0xb71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291"

func TestWhenLoadingSignedResponseSettings(t *testing.T) {
	t.Parallel()

	t.Run("Given private_key and ttl_seconds are provided", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"signed_response.private_key": testPrivateKey,
			"signed_response.ttl_seconds": 600,
		}, "."), nil)

		t.Run("Should load both and expose ttl as a duration", func(t *testing.T) {
			t.Parallel()

			result := gateway_settings.NewSignedResponseSettings(zap.NewNop(), k)

			assert.Equal(t, testPrivateKey, result.PrivateKey)
			assert.Equal(t, 600, result.TtlSeconds)
			assert.Equal(t, 10*time.Minute, result.Ttl())
		})
	})

	t.Run("Given private_key has surrounding whitespace", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"signed_response.private_key": "  " + testPrivateKey + "  ",
		}, "."), nil)

		t.Run("Should trim it", func(t *testing.T) {
			t.Parallel()

			result := gateway_settings.NewSignedResponseSettings(zap.NewNop(), k)

			assert.Equal(t, testPrivateKey, result.PrivateKey)
		})
	})

	t.Run("Given ttl_seconds is missing", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"signed_response.private_key": testPrivateKey,
		}, "."), nil)

		t.Run("Should default ttl to 300 seconds", func(t *testing.T) {
			t.Parallel()

			result := gateway_settings.NewSignedResponseSettings(zap.NewNop(), k)

			assert.Equal(t, 300, result.TtlSeconds)
			assert.Equal(t, 5*time.Minute, result.Ttl())
		})
	})

	t.Run("Given ttl_seconds is zero or negative", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"signed_response.private_key": testPrivateKey,
			"signed_response.ttl_seconds": -1,
		}, "."), nil)

		t.Run("Should fall back to the 300 second default", func(t *testing.T) {
			t.Parallel()

			result := gateway_settings.NewSignedResponseSettings(zap.NewNop(), k)

			assert.Equal(t, 300, result.TtlSeconds)
		})
	})
}

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

func TestWhenLoadingCacheSettings(t *testing.T) {
	t.Parallel()

	t.Run("Given ttl_seconds is provided", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"cache.ttl_seconds": 600,
		}, "."), nil)

		t.Run("Should load it and expose ttl as a duration", func(t *testing.T) {
			t.Parallel()

			result := gateway_settings.NewCacheSettings(zap.NewNop(), k)

			assert.Equal(t, 600, result.TtlSeconds)
			assert.Equal(t, 10*time.Minute, result.Ttl())
		})
	})

	t.Run("Given ttl_seconds is missing", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")

		t.Run("Should default ttl to 3600 seconds", func(t *testing.T) {
			t.Parallel()

			result := gateway_settings.NewCacheSettings(zap.NewNop(), k)

			assert.Equal(t, 3600, result.TtlSeconds)
			assert.Equal(t, time.Hour, result.Ttl())
		})
	})

	t.Run("Given ttl_seconds is zero or negative", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"cache.ttl_seconds": -1,
		}, "."), nil)

		t.Run("Should fall back to the 3600 second default", func(t *testing.T) {
			t.Parallel()

			result := gateway_settings.NewCacheSettings(zap.NewNop(), k)

			assert.Equal(t, 3600, result.TtlSeconds)
		})
	})
}

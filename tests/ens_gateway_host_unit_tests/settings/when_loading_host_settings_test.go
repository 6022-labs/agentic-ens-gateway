package settings_test

import (
	"testing"

	host_settings "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_host/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingHostSettings(t *testing.T) {
	t.Parallel()

	t.Run("Given a complete configuration", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"host.api_port": 8080,
		}, "."), nil)

		t.Run("Should load every field as provided", func(t *testing.T) {
			t.Parallel()

			result := host_settings.NewHostSettings(zap.NewNop(), k)

			assert.Equal(t, 8080, result.ApiPort)
		})
	})

	t.Run("Given api_port is missing", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{}, "."), nil)

		t.Run("Should default to 5000", func(t *testing.T) {
			t.Parallel()

			result := host_settings.NewHostSettings(zap.NewNop(), k)

			assert.Equal(t, 5000, result.ApiPort)
		})
	})
}

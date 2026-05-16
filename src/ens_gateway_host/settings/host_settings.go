package settings

import (
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const HostSettingsKey = "host"

type HostSettings struct {
	ApiPort int `koanf:"api_port"`
}

func NewHostSettings(logger *zap.Logger, k *koanf.Koanf) *HostSettings {
	settings := HostSettings{
		ApiPort: 5000,
	}
	if err := k.Unmarshal(HostSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal host settings", zap.Error(err))
	}

	if settings.ApiPort <= 0 {
		logger.Warn("invalid host.api_port value, defaulting to 5000")
		settings.ApiPort = 5000
	}

	return &settings
}

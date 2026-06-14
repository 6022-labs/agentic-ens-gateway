package settings

import (
	"time"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const CacheSettingsKey = "cache"

const defaultCacheTtlSeconds = 3600

type CacheSettings struct {
	TtlSeconds int `koanf:"ttl_seconds"`
}

func NewCacheSettings(logger *zap.Logger, k *koanf.Koanf) *CacheSettings {
	settings := CacheSettings{
		TtlSeconds: defaultCacheTtlSeconds,
	}
	if err := k.Unmarshal(CacheSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal cache settings", zap.Error(err))
	}

	if settings.TtlSeconds <= 0 {
		logger.Warn("invalid cache.ttl_seconds value, defaulting to 3600")
		settings.TtlSeconds = defaultCacheTtlSeconds
	}

	return &settings
}

func (s *CacheSettings) Ttl() time.Duration {
	return time.Duration(s.TtlSeconds) * time.Second
}

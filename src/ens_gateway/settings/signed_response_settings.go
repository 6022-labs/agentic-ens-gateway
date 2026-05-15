package settings

import (
	"strings"
	"time"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const SignedResponseSettingsKey = "signed_response"

type SignedResponseSettings struct {
	PrivateKey string `koanf:"private_key"`
	TtlSeconds int    `koanf:"ttl_seconds"`
}

func NewSignedResponseSettings(logger *zap.Logger, k *koanf.Koanf) *SignedResponseSettings {
	settings := SignedResponseSettings{
		TtlSeconds: 300,
	}
	if err := k.Unmarshal(SignedResponseSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal signed_response settings", zap.Error(err))
	}

	settings.PrivateKey = strings.TrimSpace(settings.PrivateKey)
	if settings.PrivateKey == "" {
		logger.Fatal("please set signed_response.private_key (or SIGNED_RESPONSE__PRIVATE_KEY env)")
	}

	if settings.TtlSeconds <= 0 {
		logger.Warn("invalid signed_response.ttl_seconds value, defaulting to 300")
		settings.TtlSeconds = 300
	}

	return &settings
}

func (s *SignedResponseSettings) Ttl() time.Duration {
	return time.Duration(s.TtlSeconds) * time.Second
}

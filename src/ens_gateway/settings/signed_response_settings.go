package settings

import (
	"strings"
	"time"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const SignedResponseSettingsKey = "signed_response"

// SignedResponseSettings controls how the gateway produces signed CCIP-Read
// responses: the ECDSA private key it signs with, and how long each signed
// response stays valid. The TTL is embedded as the `expires` field of the
// signed payload and the L1 resolver's `SignatureVerifier` rejects responses
// whose `expires` is in the past.
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
		logger.Fatal("please set signed_response.private_key (or SIGNED_RESPONSE__PRIVATE_KEY env) — 0x-prefixed hex ECDSA key the gateway signs responses with")
	}

	if settings.TtlSeconds <= 0 {
		logger.Warn("invalid signed_response.ttl_seconds value, defaulting to 300")
		settings.TtlSeconds = 300
	}

	return &settings
}

// Ttl returns the response validity window as a duration.
func (s *SignedResponseSettings) Ttl() time.Duration {
	return time.Duration(s.TtlSeconds) * time.Second
}

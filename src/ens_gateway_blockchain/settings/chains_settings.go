package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const ChainsSettingsKey = "chains"

const defaultRequestTimeoutSeconds = 10

type ChainConfig struct {
	RpcHttpUrl            string `koanf:"rpc_http_url"`
	RequestTimeoutSeconds uint64 `koanf:"request_timeout_seconds"`
}

type ChainsSettings map[uint64]*ChainConfig

func NewChainsSettings(logger *zap.Logger, k *koanf.Koanf) *ChainsSettings {
	settings := ChainsSettings{}
	if err := k.UnmarshalWithConf(ChainsSettingsKey, &settings, koanf.UnmarshalConf{Tag: "koanf"}); err != nil {
		logger.Fatal("failed to unmarshal chains settings", zap.Error(err))
	}

	if len(settings) == 0 {
		logger.Fatal(
			"chains settings is empty: configure at least one chain via appsettings.json or CHAINS__<chain_id>__* env vars",
		)
	}

	for chainId, chain := range settings {
		if chain == nil {
			logger.Fatal("chain entry is nil", zap.Uint64("chainId", chainId))
		}
		chain.RpcHttpUrl = strings.TrimSpace(chain.RpcHttpUrl)

		if chain.RpcHttpUrl == "" {
			logger.Fatal("chain rpc_http_url is required", zap.Uint64("chainId", chainId))
		}

		// An unbounded RPC call hangs the caller forever; the reference bounds every chain.
		if chain.RequestTimeoutSeconds == 0 {
			chain.RequestTimeoutSeconds = defaultRequestTimeoutSeconds
		}
	}

	return &settings
}

func (s *ChainsSettings) Get(chainId uint64) *ChainConfig {
	if s == nil {
		return nil
	}
	return (*s)[chainId]
}

package settings

import (
	_ "embed"
	"encoding/json"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

//go:embed agent_ens_registries.json
var agentEnsRegistriesJson []byte

type agentEnsRegistryEntry struct {
	ChainId uint64 `json:"chainId"`
	Address string `json:"address"`
}

type agentEnsRegistriesFile struct {
	Registries []agentEnsRegistryEntry `json:"registries"`
}

// AgentEnsRegistriesSettings holds the AgentEnsRegistry contract address per
// chain, resolved from the built-in (code-managed) registry file. This is not
// operator-facing config — it ships with the build. Empty entries are skipped
// (treated as 'not yet deployed on this chain').
type AgentEnsRegistriesSettings map[uint64]common.Address

func NewAgentEnsRegistriesSettings(logger *zap.Logger) *AgentEnsRegistriesSettings {
	var f agentEnsRegistriesFile
	if err := json.Unmarshal(agentEnsRegistriesJson, &f); err != nil {
		logger.Fatal("failed to parse built-in agent_ens_registries.json", zap.Error(err))
	}

	settings := AgentEnsRegistriesSettings{}
	for _, r := range f.Registries {
		addr := strings.TrimSpace(r.Address)
		if addr == "" {
			continue
		}
		if !common.IsHexAddress(addr) {
			logger.Fatal(
				"agent_ens_registry address is not a valid hex address",
				zap.Uint64("chainId", r.ChainId),
				zap.String("address", addr),
			)
		}
		settings[r.ChainId] = common.HexToAddress(addr)
	}

	return &settings
}

// Get returns the registry address for `chainId` and whether it was configured.
func (s *AgentEnsRegistriesSettings) Get(chainId uint64) (common.Address, bool) {
	if s == nil {
		return common.Address{}, false
	}
	addr, ok := (*s)[chainId]
	return addr, ok
}

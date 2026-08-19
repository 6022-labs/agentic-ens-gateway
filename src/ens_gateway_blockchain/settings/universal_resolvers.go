package settings

import (
	_ "embed"
	"encoding/json"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

//go:embed universal_resolvers.json
var universalResolversJson []byte

type universalResolverEntry struct {
	ChainId uint64 `json:"chainId"`
	Address string `json:"address"`
}

type sourceToResolverMappingEntry struct {
	SourceChainId            uint64 `json:"sourceChainId"`
	UniversalResolverChainId uint64 `json:"universalResolverChainId"`
}

type universalResolversFile struct {
	Resolvers []universalResolverEntry       `json:"resolvers"`
	Mappings  []sourceToResolverMappingEntry `json:"mappings"`
}

type UniversalResolverSettings struct {
	Resolvers map[uint64]common.Address
	Mappings  map[uint64]uint64
}

func NewUniversalResolverSettings(logger *zap.Logger) *UniversalResolverSettings {
	var f universalResolversFile
	if err := json.Unmarshal(universalResolversJson, &f); err != nil {
		logger.Fatal("failed to parse built-in universal_resolvers.json", zap.Error(err))
	}

	settings := &UniversalResolverSettings{
		Resolvers: map[uint64]common.Address{},
		Mappings:  map[uint64]uint64{},
	}

	for _, r := range f.Resolvers {
		addr := strings.TrimSpace(r.Address)
		if addr == "" {
			continue
		}
		if !common.IsHexAddress(addr) {
			logger.Fatal(
				"universal_resolver address is not a valid hex address",
				zap.Uint64("chainId", r.ChainId),
				zap.String("address", addr),
			)
		}
		settings.Resolvers[r.ChainId] = common.HexToAddress(addr)
	}

	for _, m := range f.Mappings {
		settings.Mappings[m.SourceChainId] = m.UniversalResolverChainId
	}

	return settings
}

// Get returns the UniversalResolver chain id and address that the given source
// chain falls back to. ok is false when no mapping exists for the source chain
// or when the mapped resolver chain has no configured address.
func (s *UniversalResolverSettings) Get(
	sourceChainId uint64,
) (universalResolverChainId uint64, addr common.Address, ok bool) {
	if s == nil {
		return 0, common.Address{}, false
	}
	universalResolverChainId, ok = s.Mappings[sourceChainId]
	if !ok {
		return 0, common.Address{}, false
	}
	addr, ok = s.Resolvers[universalResolverChainId]
	if !ok {
		return universalResolverChainId, common.Address{}, false
	}
	return universalResolverChainId, addr, true
}

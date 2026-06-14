package services

import (
	"fmt"

	ens_gateway_interfaces "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/services/interfaces"
	ens_gateway_settings "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/settings"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

// ResolutionCache caches resolved ENS values on top of the atomic cache
// backend. It owns the key layout and value encoding; backend errors are
// logged and swallowed.
type ResolutionCache struct {
	cacheStorage  ens_gateway_interfaces.CacheStorageInterface
	cacheSettings *ens_gateway_settings.CacheSettings
	logger        *zap.Logger
}

func NewResolutionCache(
	cacheStorage ens_gateway_interfaces.CacheStorageInterface,
	cacheSettings *ens_gateway_settings.CacheSettings,
	logger *zap.Logger,
) *ResolutionCache {
	return &ResolutionCache{
		cacheStorage:  cacheStorage,
		cacheSettings: cacheSettings,
		logger:        logger,
	}
}

func (c *ResolutionCache) GetText(chainId uint64, node [32]byte, key string) *string {
	return c.read(c.textKey(chainId, node, key))
}

func (c *ResolutionCache) SetText(chainId uint64, node [32]byte, key string, value string) {
	c.store(c.textKey(chainId, node, key), value)
}

func (c *ResolutionCache) GetAddr(chainId uint64, node [32]byte) *common.Address {
	cached := c.read(c.addrKey(chainId, node))
	if cached == nil {
		return nil
	}
	addr := common.HexToAddress(*cached)
	return &addr
}

func (c *ResolutionCache) SetAddr(chainId uint64, node [32]byte, value common.Address) {
	c.store(c.addrKey(chainId, node), value.Hex())
}

func (c *ResolutionCache) read(cacheKey string) *string {
	cached, err := c.cacheStorage.ReadString(cacheKey)
	if err != nil {
		c.logger.Warn("cache read failed", zap.String("key", cacheKey), zap.Error(err))
		return nil
	}
	return cached
}

func (c *ResolutionCache) store(cacheKey, value string) {
	ttl := c.cacheSettings.Ttl()
	if err := c.cacheStorage.Store(cacheKey, value, &ttl); err != nil {
		c.logger.Warn("cache store failed", zap.String("key", cacheKey), zap.Error(err))
	}
}

func (c *ResolutionCache) textKey(chainId uint64, node [32]byte, key string) string {
	return fmt.Sprintf("ens:text:%d:%s:%s", chainId, common.Bytes2Hex(node[:]), key)
}

func (c *ResolutionCache) addrKey(chainId uint64, node [32]byte) string {
	return fmt.Sprintf("ens:addr:%d:%s", chainId, common.Bytes2Hex(node[:]))
}

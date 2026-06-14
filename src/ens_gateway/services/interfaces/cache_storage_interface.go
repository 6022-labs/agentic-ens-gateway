package interfaces

import "time"

// CacheStorageInterface is the atomic cache port implemented by the Redis
// adapter in ens_gateway_redis.
type CacheStorageInterface interface {
	Store(key string, value any, expiration *time.Duration) error
	ReadString(key string) (*string, error)
}

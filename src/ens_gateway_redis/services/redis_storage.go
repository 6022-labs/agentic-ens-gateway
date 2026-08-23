package services

import (
	"context"
	"errors"
	"time"

	metrics_interfaces "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_redis/metrics/interfaces"
	"github.com/redis/go-redis/v9"
)

// RedisStorage is a thin, atomic adapter over the go-redis client.
type RedisStorage struct {
	redisClient  *redis.Client
	redisMetrics metrics_interfaces.RedisMetricsInterface
}

func NewRedisStorage(
	redisClient *redis.Client,
	redisMetrics metrics_interfaces.RedisMetricsInterface,
) *RedisStorage {
	return &RedisStorage{redisClient: redisClient, redisMetrics: redisMetrics}
}

func (rs *RedisStorage) Store(key string, value any, expiration *time.Duration) error {
	if expiration == nil {
		defaultValue := time.Duration(0)
		expiration = &defaultValue
	}

	start := time.Now()
	err := rs.redisClient.Set(context.Background(), key, value, *expiration).Err()
	rs.redisMetrics.RecordOperation(context.Background(), "SET", time.Since(start), err)
	return err
}

func (rs *RedisStorage) ReadString(key string) (*string, error) {
	start := time.Now()
	val, err := rs.redisClient.Get(context.Background(), key).Result()
	// redis.Nil is a miss, not an operational failure.
	rs.redisMetrics.RecordOperation(context.Background(), "GET", time.Since(start), rs.operationError(err))
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	return &val, nil
}

// operationError treats a missing key as a miss rather than an operational failure.
func (rs *RedisStorage) operationError(err error) error {
	if err == nil || errors.Is(err, redis.Nil) {
		return nil
	}
	return err
}

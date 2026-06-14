package services

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisStorage is a thin, atomic adapter over the go-redis client.
type RedisStorage struct {
	redisClient *redis.Client
}

func NewRedisStorage(redisClient *redis.Client) *RedisStorage {
	return &RedisStorage{redisClient: redisClient}
}

func (rs *RedisStorage) Store(key string, value any, expiration *time.Duration) error {
	if expiration == nil {
		defaultValue := time.Duration(0)
		expiration = &defaultValue
	}

	return rs.redisClient.Set(context.Background(), key, value, *expiration).Err()
}

func (rs *RedisStorage) ReadString(key string) (*string, error) {
	val, err := rs.redisClient.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	return &val, nil
}

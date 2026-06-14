package services_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_redis/services"
	"github.com/6022-labs/agentic-ens-gateway/tests/common_tests/container_builders"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

type WhenStoringValueTestingSuite struct {
	redisClient *redis.Client
	sut         *services.RedisStorage
}

func WhenStoringValueBeforeEach() *WhenStoringValueTestingSuite {
	redisClient := container_builders.GetRedisConnection()

	sut := services.NewRedisStorage(
		redisClient,
	)

	return &WhenStoringValueTestingSuite{
		redisClient: redisClient,
		sut:         sut,
	}
}

func WhenStoringValueAfterEach(suite *WhenStoringValueTestingSuite) {
	container_builders.CloseRedisConnection(suite.redisClient)
}

func TestWhenStoringValue(t *testing.T) {
	t.Parallel()

	t.Run("Given a value stored without expiration", func(t *testing.T) {
		suite := WhenStoringValueBeforeEach()
		defer WhenStoringValueAfterEach(suite)

		key := fmt.Sprintf("test_store_value_%s_%d", t.Name(), time.Now().UnixNano())
		err := suite.sut.Store(key, "value", nil)
		assert.NoError(t, err)

		t.Run("Should read the value back", func(t *testing.T) {
			result, err := suite.sut.ReadString(key)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, "value", *result)
		})
	})

	t.Run("Given a value stored with a TTL", func(t *testing.T) {
		suite := WhenStoringValueBeforeEach()
		defer WhenStoringValueAfterEach(suite)

		key := fmt.Sprintf("test_store_ttl_%s_%d", t.Name(), time.Now().UnixNano())
		ttl := time.Hour
		err := suite.sut.Store(key, "value", &ttl)
		assert.NoError(t, err)

		t.Run("Should set the expiry and read the value back before it expires", func(t *testing.T) {
			remaining := suite.redisClient.TTL(context.Background(), key).Val()
			assert.Greater(t, remaining, time.Duration(0))
			assert.LessOrEqual(t, remaining, ttl)

			result, err := suite.sut.ReadString(key)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, "value", *result)
		})
	})

	t.Run("Given a value stored with a short TTL that has elapsed", func(t *testing.T) {
		suite := WhenStoringValueBeforeEach()
		defer WhenStoringValueAfterEach(suite)

		key := fmt.Sprintf("test_store_expiry_%s_%d", t.Name(), time.Now().UnixNano())
		ttl := 50 * time.Millisecond
		err := suite.sut.Store(key, "value", &ttl)
		assert.NoError(t, err)

		t.Run("Should return nothing once the key has expired", func(t *testing.T) {
			time.Sleep(150 * time.Millisecond)

			result, err := suite.sut.ReadString(key)
			assert.NoError(t, err)
			assert.Nil(t, result)
		})
	})
}

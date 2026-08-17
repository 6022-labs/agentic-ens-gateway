package services_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_redis/services"
	"github.com/6022-labs/agentic-ens-gateway/tests/common_tests/container_builders"
	metrics_mocks "github.com/6022-labs/agentic-ens-gateway/tests/ens_gateway_redis_mocks/metrics_mocks/interfaces_mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenReadingStringTestingSuite struct {
	redisClient *redis.Client
	sut         *services.RedisStorage
}

func WhenReadingStringBeforeEach(t *testing.T) *WhenReadingStringTestingSuite {
	redisClient := container_builders.GetRedisConnection()

	redisMetrics := metrics_mocks.NewMockRedisMetricsInterface(gomock.NewController(t))
	redisMetrics.EXPECT().RecordOperation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	sut := services.NewRedisStorage(
		redisClient,
		redisMetrics,
	)

	return &WhenReadingStringTestingSuite{
		redisClient: redisClient,
		sut:         sut,
	}
}

func WhenReadingStringAfterEach(suite *WhenReadingStringTestingSuite) {
	container_builders.CloseRedisConnection(suite.redisClient)
}

func TestWhenReadingString(t *testing.T) {
	t.Parallel()

	t.Run("Given a key storing 'value' as a string", func(t *testing.T) {
		suite := WhenReadingStringBeforeEach(t)
		defer WhenReadingStringAfterEach(suite)

		key := fmt.Sprintf("test_key_value_%s_%d", t.Name(), time.Now().UnixNano())
		err := suite.redisClient.Set(context.Background(), key, "value", 0).Err()
		if err != nil {
			t.Fatalf("Failed to set key in Redis: %v", err)
		}

		t.Run("Should return the stored value", func(t *testing.T) {
			result, err := suite.sut.ReadString(key)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, "value", *result)
		})
	})

	t.Run("Given a non-existing key", func(t *testing.T) {
		suite := WhenReadingStringBeforeEach(t)
		defer WhenReadingStringAfterEach(suite)

		key := fmt.Sprintf("non_existing_%s_%d", t.Name(), time.Now().UnixNano())

		t.Run("Should return nothing without erroring", func(t *testing.T) {
			result, err := suite.sut.ReadString(key)

			assert.Nil(t, err)
			assert.Nil(t, result)
		})
	})
}

package services_test

import (
	"errors"
	"testing"
	"time"

	ens_gateway_services "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/services"
	ens_gateway_settings "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/settings"
	interfaces_mocks "github.com/6022-labs/agentic-ens-gateway/tests/ens_gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenCachingResolutionTestingSuite struct {
	mockStorage *interfaces_mocks.MockCacheStorageInterface
	sut         *ens_gateway_services.ResolutionCache
}

func whenCachingResolutionBeforeEach(t *testing.T) *WhenCachingResolutionTestingSuite {
	ctrl := gomock.NewController(t)
	mockStorage := interfaces_mocks.NewMockCacheStorageInterface(ctrl)
	settings := &ens_gateway_settings.CacheSettings{TtlSeconds: 3600}

	return &WhenCachingResolutionTestingSuite{
		mockStorage: mockStorage,
		sut:         ens_gateway_services.NewResolutionCache(mockStorage, settings, zap.NewNop()),
	}
}

func TestWhenCachingResolution(t *testing.T) {
	t.Parallel()

	var node [32]byte
	node[31] = 0x07
	const (
		chainId uint64 = 137
		key            = "url"
	)
	textCacheKey := "ens:text:137:" + common.Bytes2Hex(node[:]) + ":url"
	addrCacheKey := "ens:addr:137:" + common.Bytes2Hex(node[:])

	t.Run("Given the text value is stored", func(t *testing.T) {
		t.Parallel()
		suite := whenCachingResolutionBeforeEach(t)

		cached := "https://agent.example"
		suite.mockStorage.EXPECT().ReadString(textCacheKey).Return(&cached, nil)

		t.Run("Should read it back under the text key", func(t *testing.T) {
			result := suite.sut.GetText(chainId, node, key)

			assert.NotNil(t, result)
			assert.Equal(t, cached, *result)
		})
	})

	t.Run("Given the text value is not stored", func(t *testing.T) {
		t.Parallel()
		suite := whenCachingResolutionBeforeEach(t)

		suite.mockStorage.EXPECT().ReadString(textCacheKey).Return(nil, nil)

		t.Run("Should return nil so the caller falls through to the RPC", func(t *testing.T) {
			assert.Nil(t, suite.sut.GetText(chainId, node, key))
		})
	})

	t.Run("Given the backend read fails", func(t *testing.T) {
		t.Parallel()
		suite := whenCachingResolutionBeforeEach(t)

		suite.mockStorage.EXPECT().ReadString(textCacheKey).Return(nil, errors.New("redis down"))

		t.Run("Should swallow the error and return nil", func(t *testing.T) {
			assert.Nil(t, suite.sut.GetText(chainId, node, key))
		})
	})

	t.Run("Given a text value to store", func(t *testing.T) {
		t.Parallel()
		suite := whenCachingResolutionBeforeEach(t)

		suite.mockStorage.EXPECT().
			Store(textCacheKey, "https://agent.example", gomock.Any()).
			DoAndReturn(func(_ string, _ any, expiration *time.Duration) error {
				assert.NotNil(t, expiration)
				assert.Equal(t, time.Hour, *expiration)
				return nil
			})

		t.Run("Should store it under the text key with the configured TTL", func(t *testing.T) {
			suite.sut.SetText(chainId, node, key, "https://agent.example")
		})
	})

	t.Run("Given the backend store fails", func(t *testing.T) {
		t.Parallel()
		suite := whenCachingResolutionBeforeEach(t)

		suite.mockStorage.EXPECT().Store(textCacheKey, gomock.Any(), gomock.Any()).Return(errors.New("redis down"))

		t.Run("Should swallow the error", func(t *testing.T) {
			assert.NotPanics(t, func() {
				suite.sut.SetText(chainId, node, key, "https://agent.example")
			})
		})
	})

	t.Run("Given the addr value is stored", func(t *testing.T) {
		t.Parallel()
		suite := whenCachingResolutionBeforeEach(t)

		addr := common.HexToAddress("0x2222222222222222222222222222222222222222")
		hex := addr.Hex()
		suite.mockStorage.EXPECT().ReadString(addrCacheKey).Return(&hex, nil)

		t.Run("Should read it back as an address under the addr key", func(t *testing.T) {
			result := suite.sut.GetAddr(chainId, node)

			assert.NotNil(t, result)
			assert.Equal(t, addr, *result)
		})
	})

	t.Run("Given an addr value to store", func(t *testing.T) {
		t.Parallel()
		suite := whenCachingResolutionBeforeEach(t)

		addr := common.HexToAddress("0x3333333333333333333333333333333333333333")
		suite.mockStorage.EXPECT().
			Store(addrCacheKey, addr.Hex(), gomock.Any()).
			Return(nil)

		t.Run("Should store its hex under the addr key", func(t *testing.T) {
			suite.sut.SetAddr(chainId, node, addr)
		})
	})
}

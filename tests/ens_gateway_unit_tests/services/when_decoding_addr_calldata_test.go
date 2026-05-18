package services_test

import (
	"testing"

	gateway_services "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/services"
	ens_gateway_interfaces "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/services/interfaces"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func addrCalldata(node [32]byte) []byte {
	bytes32Type, _ := abi.NewType("bytes32", "", nil)
	packed, _ := abi.Arguments{{Type: bytes32Type}}.Pack(node)
	return append(crypto.Keccak256([]byte("addr(bytes32)"))[:4], packed...)
}

func TestWhenDecodingAddrCalldata(t *testing.T) {
	t.Parallel()

	var node [32]byte
	copy(node[:], crypto.Keccak256([]byte("test-node")))

	// `{agent}.{collection}.{chainId}.protocol6022.eth` — chain id at label index 2.
	validName := dnsEncode("myagent", "test-collection", "137", "protocol6022", "eth")

	t.Run("Given valid resolve+addr call data with a chain-id label", func(t *testing.T) {
		t.Parallel()

		t.Run("Should decode the chain id, node, dns name and inner data with QueryKindAddr", func(t *testing.T) {
			t.Parallel()

			sut := gateway_services.NewResolverServiceCodec()
			inner := addrCalldata(node)
			callData := resolveCalldata(validName, inner)

			got, err := sut.Decode(callData)

			assert.NoError(t, err)
			assert.Equal(t, ens_gateway_interfaces.QueryKindAddr, got.Kind)
			assert.Equal(t, uint64(137), got.ChainId)
			assert.Equal(t, node, got.Node)
			assert.Equal(t, "", got.Key)
			assert.Equal(t, validName, got.DnsName)
			assert.Equal(t, inner, got.InnerData)
		})
	})

	t.Run("Given an address value to encode", func(t *testing.T) {
		t.Parallel()

		t.Run("Should round-trip through abi address encoding", func(t *testing.T) {
			t.Parallel()

			sut := gateway_services.NewResolverServiceCodec()
			addr := common.HexToAddress("0x1111111111111111111111111111111111111111")

			encoded, err := sut.EncodeAddrResult(addr)
			assert.NoError(t, err)

			decoded, err := sut.DecodeAddrResult(encoded)
			assert.NoError(t, err)
			assert.Equal(t, addr, decoded)
		})
	})

	t.Run("Given an encoded address result", func(t *testing.T) {
		t.Parallel()

		t.Run("Should treat empty bytes as the zero address without error", func(t *testing.T) {
			t.Parallel()

			sut := gateway_services.NewResolverServiceCodec()

			decoded, err := sut.DecodeAddrResult(nil)
			assert.NoError(t, err)
			assert.Equal(t, common.Address{}, decoded)
		})

		t.Run("Should return an error on malformed bytes", func(t *testing.T) {
			t.Parallel()

			sut := gateway_services.NewResolverServiceCodec()

			_, err := sut.DecodeAddrResult([]byte{0xff, 0xff, 0xff})
			assert.Error(t, err)
		})
	})
}

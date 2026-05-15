package services_test

import (
	"testing"

	gateway_services "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/services"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func dnsEncode(labels ...string) []byte {
	var out []byte
	for _, l := range labels {
		out = append(out, byte(len(l)))
		out = append(out, []byte(l)...)
	}
	return append(out, 0x00)
}

func textCalldata(node [32]byte, key string) []byte {
	bytes32Type, _ := abi.NewType("bytes32", "", nil)
	stringType, _ := abi.NewType("string", "", nil)
	packed, _ := abi.Arguments{{Type: bytes32Type}, {Type: stringType}}.Pack(node, key)
	return append(crypto.Keccak256([]byte("text(bytes32,string)"))[:4], packed...)
}

func resolveCalldata(name, data []byte) []byte {
	bytesType, _ := abi.NewType("bytes", "", nil)
	packed, _ := abi.Arguments{{Type: bytesType}, {Type: bytesType}}.Pack(name, data)
	return append(crypto.Keccak256([]byte("resolve(bytes,bytes)"))[:4], packed...)
}

func TestWhenDecodingResolverCalldata(t *testing.T) {
	t.Parallel()

	var node [32]byte
	copy(node[:], crypto.Keccak256([]byte("test-node")))
	const key = "agent.runtime-url"

	// `{agent}.{collection}.{chainId}.protocol6022.eth` — chain id at label index 2.
	validName := dnsEncode("myagent", "test-collection", "137", "protocol6022", "eth")

	t.Run("Given valid resolve+text call data with a chain-id label", func(t *testing.T) {
		t.Parallel()

		t.Run("Should decode the chain id, node and key", func(t *testing.T) {
			t.Parallel()

			sut := gateway_services.NewResolverServiceCodec()
			callData := resolveCalldata(validName, textCalldata(node, key))

			gotChainId, gotNode, gotKey, err := sut.Decode(callData)

			assert.NoError(t, err)
			assert.Equal(t, uint64(137), gotChainId)
			assert.Equal(t, node, gotNode)
			assert.Equal(t, key, gotKey)
		})
	})

	t.Run("Given call data with a wrong outer selector", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			sut := gateway_services.NewResolverServiceCodec()

			_, _, _, err := sut.Decode([]byte{0xde, 0xad, 0xbe, 0xef})

			assert.Error(t, err)
		})
	})

	t.Run("Given an inner payload with a wrong selector", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			sut := gateway_services.NewResolverServiceCodec()
			callData := resolveCalldata(validName, []byte{0x11, 0x22, 0x33, 0x44})

			_, _, _, err := sut.Decode(callData)

			assert.Error(t, err)
		})
	})

	t.Run("Given truncated call data", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			sut := gateway_services.NewResolverServiceCodec()

			_, _, _, err := sut.Decode([]byte{0x01, 0x02})

			assert.Error(t, err)
		})
	})

	t.Run("Given a DNS-encoded name with fewer than 3 labels", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			sut := gateway_services.NewResolverServiceCodec()
			shortName := dnsEncode("only", "two")
			callData := resolveCalldata(shortName, textCalldata(node, key))

			_, _, _, err := sut.Decode(callData)

			assert.Error(t, err)
		})
	})

	t.Run("Given the chain-id label is not numeric", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			sut := gateway_services.NewResolverServiceCodec()
			badName := dnsEncode("myagent", "test-collection", "notanumber", "protocol6022", "eth")
			callData := resolveCalldata(badName, textCalldata(node, key))

			_, _, _, err := sut.Decode(callData)

			assert.Error(t, err)
		})
	})

	t.Run("Given a text value to encode", func(t *testing.T) {
		t.Parallel()

		t.Run("Should round-trip through abi string encoding", func(t *testing.T) {
			t.Parallel()

			sut := gateway_services.NewResolverServiceCodec()

			encoded, err := sut.EncodeResult("https://agent.example/")
			assert.NoError(t, err)

			stringType, _ := abi.NewType("string", "", nil)
			unpacked, err := abi.Arguments{{Type: stringType}}.Unpack(encoded)
			assert.NoError(t, err)
			assert.Equal(t, "https://agent.example/", unpacked[0])
		})
	})
}

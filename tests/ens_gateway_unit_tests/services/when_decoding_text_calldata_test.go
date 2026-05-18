package services_test

import (
	"testing"

	gateway_services "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/services"
	ens_gateway_interfaces "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/services/interfaces"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestWhenDecodingTextCalldata(t *testing.T) {
	t.Parallel()

	var node [32]byte
	copy(node[:], crypto.Keccak256([]byte("test-node")))
	const key = "url"

	// `{agent}.{collection}.{chainId}.protocol6022.eth` — chain id at label index 2.
	validName := dnsEncode("myagent", "test-collection", "137", "protocol6022", "eth")

	t.Run("Given valid resolve+text call data with a chain-id label", func(t *testing.T) {
		t.Parallel()

		t.Run("Should decode the chain id, node, key, dns name and inner data", func(t *testing.T) {
			t.Parallel()

			sut := gateway_services.NewResolverServiceCodec()
			inner := textCalldata(node, key)
			callData := resolveCalldata(validName, inner)

			got, err := sut.Decode(callData)

			assert.NoError(t, err)
			assert.Equal(t, ens_gateway_interfaces.QueryKindText, got.Kind)
			assert.Equal(t, uint64(137), got.ChainId)
			assert.Equal(t, node, got.Node)
			assert.Equal(t, key, got.Key)
			assert.Equal(t, validName, got.DnsName)
			assert.Equal(t, inner, got.InnerData)
		})
	})

	t.Run("Given call data with a wrong outer selector", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			sut := gateway_services.NewResolverServiceCodec()

			_, err := sut.Decode([]byte{0xde, 0xad, 0xbe, 0xef})

			assert.Error(t, err)
		})
	})

	t.Run("Given an inner payload with a wrong selector", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			sut := gateway_services.NewResolverServiceCodec()
			callData := resolveCalldata(validName, []byte{0x11, 0x22, 0x33, 0x44})

			_, err := sut.Decode(callData)

			assert.Error(t, err)
		})
	})

	t.Run("Given truncated call data", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			sut := gateway_services.NewResolverServiceCodec()

			_, err := sut.Decode([]byte{0x01, 0x02})

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

			_, err := sut.Decode(callData)

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

			_, err := sut.Decode(callData)

			assert.Error(t, err)
		})
	})

	t.Run("Given a text value to encode", func(t *testing.T) {
		t.Parallel()

		t.Run("Should round-trip through abi string encoding", func(t *testing.T) {
			t.Parallel()

			sut := gateway_services.NewResolverServiceCodec()

			encoded, err := sut.EncodeTextResult("https://agent.example/")
			assert.NoError(t, err)

			stringType, _ := abi.NewType("string", "", nil)
			unpacked, err := abi.Arguments{{Type: stringType}}.Unpack(encoded)
			assert.NoError(t, err)
			assert.Equal(t, "https://agent.example/", unpacked[0])
		})
	})

	t.Run("Given an encoded text result", func(t *testing.T) {
		t.Parallel()

		t.Run("Should decode back to the same string", func(t *testing.T) {
			t.Parallel()

			sut := gateway_services.NewResolverServiceCodec()

			encoded, err := sut.EncodeTextResult("https://agent.example/")
			assert.NoError(t, err)

			decoded, err := sut.DecodeTextResult(encoded)
			assert.NoError(t, err)
			assert.Equal(t, "https://agent.example/", decoded)
		})

		t.Run("Should treat empty bytes as empty string without error", func(t *testing.T) {
			t.Parallel()

			sut := gateway_services.NewResolverServiceCodec()

			decoded, err := sut.DecodeTextResult(nil)
			assert.NoError(t, err)
			assert.Equal(t, "", decoded)
		})

		t.Run("Should return an error on malformed bytes", func(t *testing.T) {
			t.Parallel()

			sut := gateway_services.NewResolverServiceCodec()

			_, err := sut.DecodeTextResult([]byte{0xff, 0xff, 0xff})
			assert.Error(t, err)
		})
	})
}

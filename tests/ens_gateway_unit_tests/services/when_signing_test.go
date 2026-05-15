package services_test

import (
	"testing"

	gateway_services "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/services"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestWhenSigning(t *testing.T) {
	t.Parallel()

	// Deterministic test signing key.
	key, _ := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	expectedAddress := crypto.PubkeyToAddress(key.PublicKey)

	// Any 32-byte digest works — Sign doesn't interpret it.
	digest := crypto.Keccak256([]byte("test-digest"))

	t.Run("Given a 32-byte digest", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a 65-byte signature with V normalized to 27 or 28", func(t *testing.T) {
			t.Parallel()

			sut := gateway_services.NewPrivateKeySigner(key)

			sig, err := sut.Sign(digest)

			assert.NoError(t, err)
			assert.Len(t, sig, 65)
			assert.Contains(t, []byte{27, 28}, sig[64])
		})

		t.Run("Should be deterministic for the same digest (RFC 6979)", func(t *testing.T) {
			t.Parallel()

			sut := gateway_services.NewPrivateKeySigner(key)

			first, err := sut.Sign(digest)
			assert.NoError(t, err)
			second, err := sut.Sign(digest)
			assert.NoError(t, err)

			assert.Equal(t, first, second)
		})

		t.Run("Should produce a signature that recovers to the signer's address", func(t *testing.T) {
			t.Parallel()

			sut := gateway_services.NewPrivateKeySigner(key)

			sig, err := sut.Sign(digest)
			assert.NoError(t, err)

			recoverable := make([]byte, 65)
			copy(recoverable, sig)
			recoverable[64] -= 27 // back to the 0/1 recovery id crypto.SigToPub expects
			pub, err := crypto.SigToPub(digest, recoverable)
			assert.NoError(t, err)
			assert.Equal(t, expectedAddress, crypto.PubkeyToAddress(*pub))
		})
	})

	t.Run("Given a freshly-constructed signer", func(t *testing.T) {
		t.Parallel()

		t.Run("Address should match the public key's Ethereum address", func(t *testing.T) {
			t.Parallel()

			sut := gateway_services.NewPrivateKeySigner(key)

			assert.Equal(t, expectedAddress, sut.Address())
		})
	})
}

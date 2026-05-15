package services_test

import (
	"context"
	"encoding/binary"
	"errors"
	"testing"
	"time"

	gateway_services "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/services"
	gateway_settings "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/settings"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

// failingSigner is an in-test SignerInterface fake whose Sign always errors.
type failingSigner struct{}

func (failingSigner) Sign([]byte) ([]byte, error) {
	return nil, errors.New("signer failed")
}

func (failingSigner) Address() common.Address {
	return common.Address{}
}

func TestWhenBuildingSignedResponse(t *testing.T) {
	t.Parallel()

	// Deterministic test signing key.
	key, _ := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	signerAddress := crypto.PubkeyToAddress(key.PublicKey)

	resolver := common.HexToAddress("0x1111111111111111111111111111111111111111")
	callData := []byte{0xaa, 0xbb, 0xcc}
	result := []byte("https://agent.example/")

	t.Run("Given a resolver, call data and result", func(t *testing.T) {
		t.Parallel()

		t.Run("Should produce a response the trusted signer can be recovered from", func(t *testing.T) {
			t.Parallel()

			signer := gateway_services.NewPrivateKeySigner(key)
			sut := gateway_services.NewSignedResponseBuilder(
				signer,
				&gateway_settings.SignedResponseSettings{TtlSeconds: 300},
			)

			before := time.Now().Unix()
			response, err := sut.Build(context.Background(), resolver, callData, result)
			after := time.Now().Unix()
			assert.NoError(t, err)

			bytesType, _ := abi.NewType("bytes", "", nil)
			uint64Type, _ := abi.NewType("uint64", "", nil)
			decoded, err := abi.Arguments{
				{Type: bytesType},
				{Type: uint64Type},
				{Type: bytesType},
			}.Unpack(response)
			assert.NoError(t, err)

			gotResult := decoded[0].([]byte)
			expires := decoded[1].(uint64)
			sig := decoded[2].([]byte)

			assert.Equal(t, result, gotResult)
			assert.GreaterOrEqual(t, int64(expires), before+300)
			assert.LessOrEqual(t, int64(expires), after+300)
			assert.Len(t, sig, 65)

			// Recompute keccak256(0x1900 || resolver || expires || keccak256(callData) || keccak256(result))
			// and recover the signer — exactly what the L1 SignatureVerifier does.
			expiresBytes := make([]byte, 8)
			binary.BigEndian.PutUint64(expiresBytes, expires)
			preimage := []byte{0x19, 0x00}
			preimage = append(preimage, resolver.Bytes()...)
			preimage = append(preimage, expiresBytes...)
			preimage = append(preimage, crypto.Keccak256(callData)...)
			preimage = append(preimage, crypto.Keccak256(result)...)
			digest := crypto.Keccak256(preimage)

			recoverable := make([]byte, 65)
			copy(recoverable, sig)
			recoverable[64] -= 27 // back to 0/1 recovery id for crypto.SigToPub
			pub, err := crypto.SigToPub(digest, recoverable)
			assert.NoError(t, err)
			assert.Equal(t, signerAddress, crypto.PubkeyToAddress(*pub))
		})
	})

	t.Run("Given the signer fails", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the wrapped error", func(t *testing.T) {
			t.Parallel()

			sut := gateway_services.NewSignedResponseBuilder(
				failingSigner{},
				&gateway_settings.SignedResponseSettings{TtlSeconds: 300},
			)

			response, err := sut.Build(context.Background(), resolver, callData, result)

			assert.Error(t, err)
			assert.Nil(t, response)
		})
	})
}

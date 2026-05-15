package services

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	ens_gateway_interfaces "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/services/interfaces"
	ens_gateway_settings "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/settings"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// SignedResponseBuilder produces ERC-3668 responses authenticated by an ECDSA
// signature, matching the ENS offchain-resolver `SignatureVerifier` format.
type SignedResponseBuilder struct {
	signer                 ens_gateway_interfaces.SignerInterface
	signedResponseSettings *ens_gateway_settings.SignedResponseSettings
	responseArgs           abi.Arguments
}

func NewSignedResponseBuilder(
	signer ens_gateway_interfaces.SignerInterface,
	signedResponseSettings *ens_gateway_settings.SignedResponseSettings,
) *SignedResponseBuilder {
	bytesType, _ := abi.NewType("bytes", "", nil)
	uint64Type, _ := abi.NewType("uint64", "", nil)

	return &SignedResponseBuilder{
		signer:                 signer,
		signedResponseSettings: signedResponseSettings,
		responseArgs: abi.Arguments{
			{Type: bytesType},  // result
			{Type: uint64Type}, // expires
			{Type: bytesType},  // sig
		},
	}
}

// Build produces abi.encode(result, expires, sig), signed over
// keccak256(0x1900 || resolver || expires || keccak256(callData) || keccak256(result)).
func (b *SignedResponseBuilder) Build(
	ctx context.Context,
	resolver common.Address,
	callData []byte,
	result []byte,
) ([]byte, error) {
	expires := uint64(time.Now().Add(b.signedResponseSettings.Ttl()).Unix())

	digest := makeSignatureHash(resolver, expires, callData, result)
	sig, err := b.signer.Sign(digest)
	if err != nil {
		return nil, fmt.Errorf("sign response: %w", err)
	}

	return b.responseArgs.Pack(result, expires, sig)
}

// makeSignatureHash mirrors contracts/libraries/SignatureVerifier.sol:
// keccak256(abi.encodePacked(hex"1900", target, expires, keccak256(request), keccak256(result))).
func makeSignatureHash(target common.Address, expires uint64, request, result []byte) []byte {
	expiresBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(expiresBytes, expires)

	preimage := make([]byte, 0, 2+common.AddressLength+8+32+32)
	preimage = append(preimage, 0x19, 0x00)
	preimage = append(preimage, target.Bytes()...)
	preimage = append(preimage, expiresBytes...)
	preimage = append(preimage, crypto.Keccak256(request)...)
	preimage = append(preimage, crypto.Keccak256(result)...)

	return crypto.Keccak256(preimage)
}

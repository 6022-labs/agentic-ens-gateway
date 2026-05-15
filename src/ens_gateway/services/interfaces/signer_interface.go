package interfaces

import "github.com/ethereum/go-ethereum/common"

// SignerInterface signs response digests for the gateway. The construction path
// (Vault-backed key, etc.) is hidden behind this port.
type SignerInterface interface {
	// Sign signs a 32-byte digest, returning a 65-byte [R || S || V] signature
	// with V normalized to 27/28.
	Sign(digest []byte) ([]byte, error)

	// Address is the signer's Ethereum address.
	Address() common.Address
}

package interfaces

// ResolverServiceCodecInterface decodes ERC-3668 resolver call data and encodes
// resolved record values back into the resolver result shape.
type ResolverServiceCodecInterface interface {
	// Decode unwraps IResolverService.resolve(bytes name, bytes data), extracts
	// the chain id from a fixed label position of the DNS-encoded name (the
	// `{agent}.{collection}.{chainId}.protocol6022.eth` convention — label index
	// 2), and decodes the inner text(bytes32 node, string key) query.
	Decode(callData []byte) (chainId uint64, node [32]byte, key string, err error)

	// EncodeResult ABI-encodes a resolved text value as the resolver result.
	EncodeResult(value string) ([]byte, error)
}

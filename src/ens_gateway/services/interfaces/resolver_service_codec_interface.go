package interfaces

type ResolverServiceCodecInterface interface {
	Decode(callData []byte) (chainId uint64, node [32]byte, key string, err error)
	EncodeResult(value string) ([]byte, error)
}

package interfaces

type DecodedQuery struct {
	ChainId   uint64
	Node      [32]byte
	Key       string
	DnsName   []byte
	InnerData []byte
}

type ResolverServiceCodecInterface interface {
	Decode(callData []byte) (DecodedQuery, error)
	EncodeResult(value string) ([]byte, error)
	DecodeTextResult(result []byte) (string, error)
}

package interfaces

import "github.com/ethereum/go-ethereum/common"

// QueryKind tells the use case which inner ENS resolver call the client made,
// so it can dispatch to the right reader method and result encoder.
type QueryKind int

const (
	QueryKindText QueryKind = iota
	QueryKindAddr
)

type DecodedQuery struct {
	Kind      QueryKind
	ChainId   uint64
	Node      [32]byte
	Key       string // populated for QueryKindText only
	DnsName   []byte
	InnerData []byte
}

type ResolverServiceCodecInterface interface {
	Decode(callData []byte) (DecodedQuery, error)
	EncodeTextResult(value string) ([]byte, error)
	DecodeTextResult(result []byte) (string, error)
	EncodeAddrResult(value common.Address) ([]byte, error)
	DecodeAddrResult(result []byte) (common.Address, error)
}

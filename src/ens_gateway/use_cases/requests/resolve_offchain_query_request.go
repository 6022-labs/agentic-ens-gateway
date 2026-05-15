package requests

import "github.com/ethereum/go-ethereum/common"

// ResolveOffchainQueryRequest is a decoded ERC-3668 CCIP-Read request: the
// resolver that issued the OffchainLookup, and the resolver call data.
type ResolveOffchainQueryRequest struct {
	Sender common.Address
	Data   []byte
}

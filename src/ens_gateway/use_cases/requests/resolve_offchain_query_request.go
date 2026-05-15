package requests

import "github.com/ethereum/go-ethereum/common"

type ResolveOffchainQueryRequest struct {
	Sender common.Address
	Data   []byte
}

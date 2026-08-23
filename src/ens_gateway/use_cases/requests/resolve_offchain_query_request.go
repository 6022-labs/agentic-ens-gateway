package requests

import (
	"github.com/6022-labs/agentic-ens-gateway/src/common/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// ResolveOffchainQueryRequest is the CCIP-Read sender plus its 0x-hex resolver call data.
type ResolveOffchainQueryRequest struct {
	Sender string `uri:"sender" json:"sender"`
	Data   string `uri:"data"   json:"data"`

	sender common.Address
	data   []byte
}

func (r *ResolveOffchainQueryRequest) ValidateAndSanitize() error {
	if !common.IsHexAddress(r.Sender) {
		return errors.NewValidationError("sender", "invalid sender address")
	}

	data, err := hexutil.Decode(r.Data)
	if err != nil {
		return errors.NewValidationError("data", "invalid call data hex")
	}

	r.sender = common.HexToAddress(r.Sender)
	r.data = data

	return nil
}

// SenderAddress is the parsed sender; only valid after ValidateAndSanitize.
func (r *ResolveOffchainQueryRequest) SenderAddress() common.Address {
	return r.sender
}

// CallData is the decoded call data; only valid after ValidateAndSanitize.
func (r *ResolveOffchainQueryRequest) CallData() []byte {
	return r.data
}

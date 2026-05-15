package interfaces

import "github.com/ethereum/go-ethereum/common"

type SignerInterface interface {
	Sign(digest []byte) ([]byte, error)
	Address() common.Address
}

package interfaces

import "github.com/ethereum/go-ethereum/common"

// ResolutionCacheInterface caches resolved ENS values. A nil return means a
// miss (or a swallowed backend error), so the caller falls through to the RPC.
type ResolutionCacheInterface interface {
	GetText(chainId uint64, node [32]byte, key string) *string
	SetText(chainId uint64, node [32]byte, key string, value string)
	GetAddr(chainId uint64, node [32]byte) *common.Address
	SetAddr(chainId uint64, node [32]byte, value common.Address)
}

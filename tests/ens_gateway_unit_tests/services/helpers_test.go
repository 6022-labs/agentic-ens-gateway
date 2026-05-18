package services_test

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
)

// dnsEncode produces a DNS-wire-format ENS name from a slice of labels —
// each label is length-prefixed and the sequence is null-terminated.
func dnsEncode(labels ...string) []byte {
	var out []byte
	for _, l := range labels {
		out = append(out, byte(len(l)))
		out = append(out, []byte(l)...)
	}
	return append(out, 0x00)
}

// textCalldata builds the inner `text(bytes32,string)` calldata that an L1
// resolver forwards through CCIP-Read.
func textCalldata(node [32]byte, key string) []byte {
	bytes32Type, _ := abi.NewType("bytes32", "", nil)
	stringType, _ := abi.NewType("string", "", nil)
	packed, _ := abi.Arguments{{Type: bytes32Type}, {Type: stringType}}.Pack(node, key)
	return append(crypto.Keccak256([]byte("text(bytes32,string)"))[:4], packed...)
}

// addrCalldata builds the inner `addr(bytes32)` calldata that an L1 resolver
// forwards through CCIP-Read.
func addrCalldata(node [32]byte) []byte {
	bytes32Type, _ := abi.NewType("bytes32", "", nil)
	packed, _ := abi.Arguments{{Type: bytes32Type}}.Pack(node)
	return append(crypto.Keccak256([]byte("addr(bytes32)"))[:4], packed...)
}

// resolveCalldata wraps an inner record query in the outer ENSIP-10
// `resolve(bytes,bytes)` envelope alongside the DNS-encoded name.
func resolveCalldata(name, data []byte) []byte {
	bytesType, _ := abi.NewType("bytes", "", nil)
	packed, _ := abi.Arguments{{Type: bytesType}, {Type: bytesType}}.Pack(name, data)
	return append(crypto.Keccak256([]byte("resolve(bytes,bytes)"))[:4], packed...)
}

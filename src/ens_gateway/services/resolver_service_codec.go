package services

import (
	"bytes"
	"fmt"
	"strconv"

	ens_gateway_interfaces "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/services/interfaces"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Selectors for the ABI layers the gateway decodes. resolve(bytes,bytes) is
// the ENSIP-10 / IResolverService outer call; text(bytes32,string) and
// addr(bytes32) are the inner ENS record queries the L1 resolver forwards.
var (
	resolveSelector = crypto.Keccak256([]byte("resolve(bytes,bytes)"))[:4]
	textSelector    = crypto.Keccak256([]byte("text(bytes32,string)"))[:4]
	addrSelector    = crypto.Keccak256([]byte("addr(bytes32)"))[:4]
)

// chainIdLabelIndex is the position of the chain-id label in the DNS-encoded
// ENS name, counting from the leftmost label. Matches the existing domain
// convention `{agent}.{collection}.{chainId}.protocol6022.eth`.
const chainIdLabelIndex = 2

// ResolverServiceCodec decodes ERC-3668 resolver call data and encodes resolved
// values back into the resolver result shape.
type ResolverServiceCodec struct {
	resolveArgs    abi.Arguments
	textArgs       abi.Arguments
	addrArgs       abi.Arguments
	textResultArgs abi.Arguments
	addrResultArgs abi.Arguments
}

func NewResolverServiceCodec() *ResolverServiceCodec {
	bytesType, _ := abi.NewType("bytes", "", nil)
	bytes32Type, _ := abi.NewType("bytes32", "", nil)
	stringType, _ := abi.NewType("string", "", nil)
	addressType, _ := abi.NewType("address", "", nil)

	return &ResolverServiceCodec{
		resolveArgs:    abi.Arguments{{Type: bytesType}, {Type: bytesType}},
		textArgs:       abi.Arguments{{Type: bytes32Type}, {Type: stringType}},
		addrArgs:       abi.Arguments{{Type: bytes32Type}},
		textResultArgs: abi.Arguments{{Type: stringType}},
		addrResultArgs: abi.Arguments{{Type: addressType}},
	}
}

// Decode unwraps the outer IResolverService.resolve(bytes name, bytes data)
// call, extracts the chain id from a fixed label of the DNS-encoded name, and
// decodes the inner query — either text(bytes32 node, string key) or
// addr(bytes32 node).
func (c *ResolverServiceCodec) Decode(callData []byte) (ens_gateway_interfaces.DecodedQuery, error) {
	if len(callData) < 4 || !bytes.Equal(callData[:4], resolveSelector) {
		return ens_gateway_interfaces.DecodedQuery{}, fmt.Errorf("unexpected call data selector, want resolve(bytes,bytes)")
	}

	outer, err := c.resolveArgs.Unpack(callData[4:])
	if err != nil {
		return ens_gateway_interfaces.DecodedQuery{}, fmt.Errorf("decode resolve(bytes,bytes): %w", err)
	}
	name, ok := outer[0].([]byte)
	if !ok {
		return ens_gateway_interfaces.DecodedQuery{}, fmt.Errorf("decode resolve(bytes,bytes): name is not bytes")
	}
	innerData, ok := outer[1].([]byte)
	if !ok {
		return ens_gateway_interfaces.DecodedQuery{}, fmt.Errorf("decode resolve(bytes,bytes): inner data is not bytes")
	}

	chainId, err := extractChainIdFromDnsName(name, chainIdLabelIndex)
	if err != nil {
		return ens_gateway_interfaces.DecodedQuery{}, fmt.Errorf("extract chain id: %w", err)
	}

	if len(innerData) < 4 {
		return ens_gateway_interfaces.DecodedQuery{}, fmt.Errorf("inner call data is truncated")
	}
	innerSelector := innerData[:4]

	switch {
	case bytes.Equal(innerSelector, textSelector):
		return c.decodeText(chainId, name, innerData)
	case bytes.Equal(innerSelector, addrSelector):
		return c.decodeAddr(chainId, name, innerData)
	default:
		return ens_gateway_interfaces.DecodedQuery{}, fmt.Errorf("unsupported inner selector, want text(bytes32,string) or addr(bytes32)")
	}
}

func (c *ResolverServiceCodec) decodeText(chainId uint64, name, innerData []byte) (ens_gateway_interfaces.DecodedQuery, error) {
	inner, err := c.textArgs.Unpack(innerData[4:])
	if err != nil {
		return ens_gateway_interfaces.DecodedQuery{}, fmt.Errorf("decode text(bytes32,string): %w", err)
	}
	node, ok := inner[0].([32]byte)
	if !ok {
		return ens_gateway_interfaces.DecodedQuery{}, fmt.Errorf("decode text(bytes32,string): node is not bytes32")
	}
	key, ok := inner[1].(string)
	if !ok {
		return ens_gateway_interfaces.DecodedQuery{}, fmt.Errorf("decode text(bytes32,string): key is not string")
	}
	return ens_gateway_interfaces.DecodedQuery{
		Kind:      ens_gateway_interfaces.QueryKindText,
		ChainId:   chainId,
		Node:      node,
		Key:       key,
		DnsName:   name,
		InnerData: innerData,
	}, nil
}

func (c *ResolverServiceCodec) decodeAddr(chainId uint64, name, innerData []byte) (ens_gateway_interfaces.DecodedQuery, error) {
	inner, err := c.addrArgs.Unpack(innerData[4:])
	if err != nil {
		return ens_gateway_interfaces.DecodedQuery{}, fmt.Errorf("decode addr(bytes32): %w", err)
	}
	node, ok := inner[0].([32]byte)
	if !ok {
		return ens_gateway_interfaces.DecodedQuery{}, fmt.Errorf("decode addr(bytes32): node is not bytes32")
	}
	return ens_gateway_interfaces.DecodedQuery{
		Kind:      ens_gateway_interfaces.QueryKindAddr,
		ChainId:   chainId,
		Node:      node,
		DnsName:   name,
		InnerData: innerData,
	}, nil
}

// EncodeTextResult ABI-encodes a resolved text value as the resolver result bytes.
func (c *ResolverServiceCodec) EncodeTextResult(value string) ([]byte, error) {
	return c.textResultArgs.Pack(value)
}

// DecodeTextResult is the inverse of EncodeTextResult.
func (c *ResolverServiceCodec) DecodeTextResult(result []byte) (string, error) {
	if len(result) == 0 {
		return "", nil
	}
	unpacked, err := c.textResultArgs.Unpack(result)
	if err != nil {
		return "", fmt.Errorf("decode text result: %w", err)
	}
	value, ok := unpacked[0].(string)
	if !ok {
		return "", fmt.Errorf("decode text result: not a string")
	}
	return value, nil
}

// EncodeAddrResult ABI-encodes a resolved address as the resolver result bytes.
func (c *ResolverServiceCodec) EncodeAddrResult(value common.Address) ([]byte, error) {
	return c.addrResultArgs.Pack(value)
}

// DecodeAddrResult is the inverse of EncodeAddrResult.
func (c *ResolverServiceCodec) DecodeAddrResult(result []byte) (common.Address, error) {
	if len(result) == 0 {
		return common.Address{}, nil
	}
	unpacked, err := c.addrResultArgs.Unpack(result)
	if err != nil {
		return common.Address{}, fmt.Errorf("decode addr result: %w", err)
	}
	value, ok := unpacked[0].(common.Address)
	if !ok {
		return common.Address{}, fmt.Errorf("decode addr result: not an address")
	}
	return value, nil
}

// extractChainIdFromDnsName parses a DNS-encoded ENS name (length-prefixed
// labels, null-terminated) and returns the label at `labelIndex` (counted from
// the leftmost label, zero-based) parsed as a uint64.
func extractChainIdFromDnsName(name []byte, labelIndex int) (uint64, error) {
	pos := 0
	for i := 0; ; i++ {
		if pos >= len(name) {
			return 0, fmt.Errorf("DNS-encoded name truncated before label %d", labelIndex)
		}
		labelLen := int(name[pos])
		if labelLen == 0 {
			return 0, fmt.Errorf("DNS-encoded name has fewer than %d labels", labelIndex+1)
		}
		pos++
		if pos+labelLen > len(name) {
			return 0, fmt.Errorf("DNS-encoded name has a label that overruns the buffer")
		}
		if i == labelIndex {
			value, err := strconv.ParseUint(string(name[pos:pos+labelLen]), 10, 64)
			if err != nil {
				return 0, fmt.Errorf("label %d is not a numeric chain id: %w", labelIndex, err)
			}
			return value, nil
		}
		pos += labelLen
	}
}

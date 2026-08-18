package model_builders

import (
	ens_gateway_interfaces "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/services/interfaces"
)

type DecodedQueryBuilder struct {
	kind      ens_gateway_interfaces.QueryKind
	chainId   uint64
	node      [32]byte
	key       string
	dnsName   []byte
	innerData []byte
}

func NewDecodedQueryBuilder() *DecodedQueryBuilder {
	return &DecodedQueryBuilder{
		kind:      ens_gateway_interfaces.QueryKindText,
		chainId:   80002,
		node:      [32]byte{0x01},
		key:       "url",
		dnsName:   []byte{0x01, 0x02},
		innerData: []byte{0x03, 0x04},
	}
}

func (b *DecodedQueryBuilder) WithKind(kind ens_gateway_interfaces.QueryKind) *DecodedQueryBuilder {
	b.kind = kind
	return b
}

func (b *DecodedQueryBuilder) WithChainId(chainId uint64) *DecodedQueryBuilder {
	b.chainId = chainId
	return b
}

func (b *DecodedQueryBuilder) WithNode(node [32]byte) *DecodedQueryBuilder {
	b.node = node
	return b
}

func (b *DecodedQueryBuilder) WithKey(key string) *DecodedQueryBuilder {
	b.key = key
	return b
}

func (b *DecodedQueryBuilder) WithDnsName(dnsName []byte) *DecodedQueryBuilder {
	b.dnsName = dnsName
	return b
}

func (b *DecodedQueryBuilder) WithInnerData(innerData []byte) *DecodedQueryBuilder {
	b.innerData = innerData
	return b
}

func (b *DecodedQueryBuilder) Build() ens_gateway_interfaces.DecodedQuery {
	return ens_gateway_interfaces.DecodedQuery{
		Kind:      b.kind,
		ChainId:   b.chainId,
		Node:      b.node,
		Key:       b.key,
		DnsName:   b.dnsName,
		InnerData: b.innerData,
	}
}

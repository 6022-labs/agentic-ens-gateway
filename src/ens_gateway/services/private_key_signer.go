package services

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type PrivateKeySigner struct {
	key *ecdsa.PrivateKey
}

func NewPrivateKeySigner(key *ecdsa.PrivateKey) *PrivateKeySigner {
	return &PrivateKeySigner{key: key}
}

// Sign returns 65-byte [R || S || V] with V normalized to 27/28 — the format OZ's ECDSA.recover expects.
func (s *PrivateKeySigner) Sign(digest []byte) ([]byte, error) {
	sig, err := crypto.Sign(digest, s.key)
	if err != nil {
		return nil, err
	}
	sig[64] += 27
	return sig, nil
}

func (s *PrivateKeySigner) Address() common.Address {
	return crypto.PubkeyToAddress(s.key.PublicKey)
}

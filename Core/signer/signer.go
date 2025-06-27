package signer

import (
	"FichainCore/common"
	"FichainCore/crypto"
	"FichainCore/types"
)

type Signer struct {
	privateKey types.PrivateKey
}

func NewSigner(
	privateKey types.PrivateKey,
) *Signer {
	return &Signer{
		privateKey: privateKey,
	}
}

func (s *Signer) SignBytes(b []byte) (types.Sign, error) {
	hash := crypto.Keccak256Hash(b)
	return s.SignHash(hash)
}

func (s *Signer) SignHash(hash common.Hash) (types.Sign, error) {
	ecdsaPK, err := crypto.ToECDSA(s.privateKey.Bytes())
	if err != nil {
		return types.Sign{}, err
	}
	bSign, err := crypto.Sign(hash.Bytes(), ecdsaPK)
	if err != nil {
		return types.Sign{}, err
	}
	sign := types.SignFromBytes(bSign)
	return sign, nil
}

func (s *Signer) WalletAddress() (common.Address, error) {
	pri, err := crypto.ToECDSA(s.privateKey.Bytes())
	if err != nil {
		return common.Address{}, nil
	}
	pub := pri.PublicKey
	address := crypto.PubkeyToAddress(pub)
	return address, nil
}

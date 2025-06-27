package types

import (
	"encoding/hex"
)

type Commands string

type (
	PrivateKey [32]byte
	PublicKey  [64]byte
	Sign       [65]byte
)

func (pk PublicKey) Bytes() []byte {
	return pk[:]
}

func (pk PublicKey) String() string {
	return hex.EncodeToString(pk[:])
}

func (pk PrivateKey) Bytes() []byte {
	return pk[:]
}

func (pk PrivateKey) String() string {
	return hex.EncodeToString(pk[:])
}

func (s Sign) Bytes() []byte {
	return s[:]
}

func (s Sign) String() string {
	return hex.EncodeToString(s[:])
}

func PubkeyFromBytes(bytes []byte) PublicKey {
	p := PublicKey{}
	copy(p[0:48], bytes)
	return p
}

func SignFromBytes(bytes []byte) Sign {
	s := Sign{}
	copy(s[0:65], bytes)
	return s
}

func PrivateKeyFromBytes(bytes []byte) PrivateKey {
	p := PrivateKey{}
	copy(p[0:32], bytes)
	return p
}

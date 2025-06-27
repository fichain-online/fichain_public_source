package transaction

import (
	"fmt"
	"log/slog"
	"math/big"
	"testing"

	"FichainCore/common"
	"github.com/stretchr/testify/assert"

	"FichainCore/crypto"
	"FichainCore/signer"
	"FichainCore/types"
)

func TestNewTransaction(t *testing.T) {
	transaction := NewTransaction(
		common.HexToAddress("431c5a61cabff377c36f34076a18a1590115f58e"),
		big.NewInt(1),
		big.NewInt(1),
		[]byte{1, 2, 3},
		100000,
		100000,
		"test message",
	)
	fmt.Printf("My transaction: %v", transaction)
	slog.Info(fmt.Sprintf("%v", transaction))
	assert.NotNil(t, transaction)
}

func TestGetTransactionHash(t *testing.T) {
	transaction := NewTransaction(
		common.HexToAddress("431c5a61cabff377c36f34076a18a1590115f58e"),
		big.NewInt(1),
		big.NewInt(1),
		[]byte{1, 2, 3},
		100000,
		100000,
		"test message",
	)
	hash := transaction.Hash()
	slog.Info(fmt.Sprintf("%v", hash))
	assert.NotNil(t, hash)
	assert.Equal(t, hash, transaction.hash)
}

func TestSignAndVerifyTransactionSign(t *testing.T) {
	transaction := NewTransaction(
		common.HexToAddress("431c5a61cabff377c36f34076a18a1590115f58e"),
		big.NewInt(1),
		big.NewInt(1),
		[]byte{1, 2, 3},
		100000,
		100000,
		"test message",
	)
	hash := transaction.Hash()
	sgn := signer.NewSigner(
		types.PrivateKeyFromBytes(
			common.FromHex("edd4932bdf38f4f632ea1fd723188d013d563feb7a5ec0c83e0fe58d3941b0fd"),
		),
	)
	sign, err := sgn.SignHash(hash)
	assert.Nil(t, err)
	signValid, err := crypto.VerifySignature(
		hash.Bytes(),
		sign,
		common.HexToAddress("0x4d5977aa43731a81c1CCB318421B3EF750CD4176"),
	)
	assert.Nil(t, err)
	assert.True(t, signValid)
}

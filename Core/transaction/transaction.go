package transaction

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/big"

	"google.golang.org/protobuf/proto"

	e_common "FichainCore/common"
	"FichainCore/crypto"
	pb "FichainCore/proto"
	"FichainCore/types"
)

type Transaction struct {
	hash      e_common.Hash
	toAddress e_common.Address
	nonce     uint64
	amount    *big.Int
	data      []byte
	gas       uint64   // changed from uint64
	gasPrice  *big.Int // changed from uint64
	message   string

	sign types.Sign
}

func NewTransaction(
	toAddress e_common.Address,
	nonce uint64,
	amount *big.Int,
	data []byte,
	gas uint64,
	gasPrice *big.Int,
	message string,
) *Transaction {
	return &Transaction{
		toAddress: toAddress,
		nonce:     nonce,
		amount:    amount,
		data:      data,
		gas:       gas,
		gasPrice:  gasPrice,
		message:   message,
	}
}

func (t *Transaction) Hash() e_common.Hash {
	if (t.hash == e_common.Hash{}) {
		var err error
		t.hash, err = t.calculateHash()
		if err != nil {
			slog.Error(fmt.Sprintf("error when calculate transaction hash %v", err))
			return e_common.Hash{}
		}
	}
	return t.hash
}

func (t *Transaction) calculateHash() (e_common.Hash, error) {
	hashData := &pb.TransactionHashData{
		ToAddress: t.toAddress.Bytes(),
		Nonce:     t.nonce,
		Amount:    t.amount.Bytes(),
		Data:      t.data,
		Gas:       t.gas,
		GasPrice:  t.gasPrice.Bytes(),
		Message:   t.message,
		Sign:      t.sign.Bytes(),
	}
	bHashData, err := proto.Marshal(hashData)
	if err != nil {
		return e_common.Hash{}, err
	}
	return crypto.Keccak256Hash(bHashData), nil
}

func (t *Transaction) HashSign(chainId *big.Int) (e_common.Hash, error) {
	signData := &pb.TransactionSignData{
		ToAddress: t.To().Bytes(),
		Nonce:     t.Nonce(),
		Amount:    t.Amount().Bytes(),
		Data:      t.Data(),
		Gas:       t.Gas(),
		GasPrice:  t.GasPrice().Bytes(),
		Message:   t.Message(),
		ChainId:   chainId.Bytes(),
	}
	b, err := proto.Marshal(signData)
	if err != nil {
		return e_common.Hash{}, err
	}
	return crypto.Keccak256Hash(b), nil
}

func (t *Transaction) From(chainId *big.Int) (e_common.Address, error) {
	if (t.sign == types.Sign{}) {
		return e_common.Address{}, fmt.Errorf("not signed yet")
	}
	hash, err := t.HashSign(chainId)
	if err != nil {
		return e_common.Address{}, err
	}

	publicKey, err := crypto.SigToPub(hash.Bytes(), t.sign.Bytes())
	if err != nil {
		return e_common.Address{}, err
	}

	return crypto.PubkeyToAddress(*publicKey), nil
}

func (t *Transaction) To() e_common.Address {
	return t.toAddress
}

func (t *Transaction) Nonce() uint64 {
	return t.nonce
}

func (t *Transaction) Amount() *big.Int {
	return t.amount
}

func (t *Transaction) Data() []byte {
	return t.data
}

func (t *Transaction) Gas() uint64 {
	return t.gas
}

func (t *Transaction) GasPrice() *big.Int {
	return t.gasPrice
}

func (t *Transaction) Message() string {
	return t.message
}

func (t *Transaction) Sign() types.Sign {
	return t.sign
}

func (t *Transaction) SetSign(sign types.Sign) {
	t.sign = sign
}

func (t *Transaction) String() string {
	return fmt.Sprintf("Transaction Hash: %s\n"+
		"To Address: %s\n"+
		"Nonce: %d\n"+
		"Amount: %s\n"+
		"Data: %s\n"+
		"Max Gas: %d\n"+
		"Max Gas Price: %s\n"+
		"Message: %s\n"+
		"Sign: %s",
		t.Hash().Hex(),
		t.To().Hex(),
		t.Nonce(),
		t.Amount(),
		hex.EncodeToString(t.Data()),
		t.Gas(),
		t.GasPrice().String(),
		t.Message(),
		t.Sign().String())
}

func (t *Transaction) Marshal() ([]byte, error) {
	return proto.Marshal(t.Proto())
}

func (t *Transaction) Unmarshal(b []byte) error {
	pbTx := &pb.Transaction{}
	err := proto.Unmarshal(b, pbTx)
	if err != nil {
		return err
	}
	t.FromProto(pbTx)
	return nil
}

func (t *Transaction) Proto() proto.Message {
	return &pb.Transaction{
		ToAddress: t.toAddress.Bytes(),
		Nonce:     t.nonce,
		Amount:    t.amount.Bytes(),
		Data:      t.data,
		Gas:       t.gas,
		GasPrice:  t.gasPrice.Bytes(),
		Message:   t.message,
		Sign:      t.sign.Bytes(),
		Hash:      t.hash.Bytes(),
	}
}

func (t *Transaction) FromProto(pbTx *pb.Transaction) error {
	t.toAddress = e_common.BytesToAddress(pbTx.ToAddress)
	t.nonce = pbTx.Nonce
	t.amount = new(big.Int).SetBytes(pbTx.Amount)
	t.data = pbTx.Data
	t.gas = pbTx.Gas
	t.gasPrice = new(big.Int).SetBytes(pbTx.GasPrice)
	t.message = pbTx.Message
	t.sign = types.SignFromBytes(pbTx.Sign)
	t.hash = e_common.BytesToHash(pbTx.Hash)
	return nil
}

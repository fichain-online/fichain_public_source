package call_data

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"FichainCore/common"
	"FichainCore/crypto"
	pb "FichainCore/proto" // make sure this matches your proto import
	"FichainCore/types"
)

type CallSmartContractData struct {
	toAddress common.Address
	data      []byte
	sign      types.Sign
}

func NewCallSmartContractData(to common.Address, data []byte) *CallSmartContractData {
	return &CallSmartContractData{
		toAddress: to,
		data:      data,
	}
}

func (c *CallSmartContractData) To() common.Address {
	return c.toAddress
}

func (c *CallSmartContractData) Data() []byte {
	return c.data
}

func (c *CallSmartContractData) Sign() types.Sign {
	return c.sign
}

func (c *CallSmartContractData) SetSign(sign types.Sign) {
	c.sign = sign
}

// For generating hash used for signing
func (c *CallSmartContractData) HashSign() (common.Hash, error) {
	msg := &pb.CallSmartContractHashData{
		ToAddress: c.toAddress.Bytes(),
		Data:      c.data,
	}
	bz, err := proto.Marshal(msg)
	if err != nil {
		return common.Hash{}, err
	}
	return crypto.Keccak256Hash(bz), nil
}

// Recover sender address from signature
func (c *CallSmartContractData) From() (common.Address, error) {
	if (c.sign == types.Sign{}) {
		return common.Address{}, fmt.Errorf("not signed yet")
	}
	hash, err := c.HashSign()
	if err != nil {
		return common.Address{}, err
	}
	pubKey, err := crypto.SigToPub(hash.Bytes(), c.sign.Bytes())
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*pubKey), nil
}

// Protobuf serialization
func (c *CallSmartContractData) Proto() proto.Message {
	return &pb.CallSmartContractData{
		ToAddress: c.toAddress.Bytes(),
		Data:      c.data,
		Sign:      c.sign.Bytes(),
	}
}

func (c *CallSmartContractData) Marshal() ([]byte, error) {
	return proto.Marshal(c.Proto())
}

func (c *CallSmartContractData) Unmarshal(bz []byte) error {
	pbData := &pb.CallSmartContractData{}
	if err := proto.Unmarshal(bz, pbData); err != nil {
		return err
	}
	return c.FromProto(pbData)
}

func (c *CallSmartContractData) FromProto(pbData *pb.CallSmartContractData) error {
	c.toAddress = common.BytesToAddress(pbData.ToAddress)
	c.data = pbData.Data
	c.sign = types.SignFromBytes(pbData.Sign)
	return nil
}

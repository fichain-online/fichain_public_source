package receipt

import (
	"fmt"
	"math/big"

	"google.golang.org/protobuf/proto"

	"FichainCore/common"
	"FichainCore/crypto"
	"FichainCore/log"
	pb "FichainCore/proto"
)

const (
	RECEIP_STATUS_SUCCESS = 1
	RECEIP_STATUS_REVERT  = 2
)

type Receipt struct {
	TxHash            common.Hash
	PostState         common.Hash
	BlockHash         common.Hash
	BlockNumber       uint64
	TxIndex           uint32
	From              common.Address
	To                common.Address // Optional; nil for contract creation
	Amount            *big.Int
	CumulativeGasUsed uint64
	GasUsed           uint64
	ContractAddress   common.Address // Optional; nil if not contract creation
	Logs              []*log.Log
	LogsBloom         []byte // 256 bytes
	Status            uint32 // 1 = success, 0 = failure
}

func NewReceipt(postState []byte, failed bool, cumulativeGasUsed uint64) *Receipt {
	return &Receipt{
		PostState:         common.BytesToHash(postState),
		Status:            boolToUint32(!failed),
		CumulativeGasUsed: cumulativeGasUsed,
		Amount:            new(big.Int),
	}
}

func boolToUint32(ok bool) uint32 {
	if ok {
		return 1
	}
	return 0
}

func (r *Receipt) Proto() proto.Message {
	var to []byte
	if (r.To != common.Address{}) {
		to = r.To.Bytes()
	}

	var contractAddr []byte
	if (r.ContractAddress != common.Address{}) {
		contractAddr = r.ContractAddress.Bytes()
	}

	logs := make([]*pb.Log, len(r.Logs))
	for i, log := range r.Logs {
		logs[i] = log.Proto()
	}

	return &pb.Receipt{
		TxHash:            r.TxHash.Bytes(),
		BlockHash:         r.BlockHash.Bytes(),
		BlockNumber:       r.BlockNumber,
		TxIndex:           r.TxIndex,
		From:              r.From.Bytes(),
		To:                to,
		Amount:            r.Amount.Bytes(),
		CumulativeGasUsed: r.CumulativeGasUsed,
		GasUsed:           r.GasUsed,
		ContractAddress:   contractAddr,
		Logs:              logs,
		PostState:         r.PostState.Bytes(),
		LogsBloom:         r.LogsBloom,
		Status:            r.Status,
	}
}

func (r *Receipt) FromProto(pbReceipt *pb.Receipt) error {
	r.TxHash = common.BytesToHash(pbReceipt.TxHash)
	r.BlockHash = common.BytesToHash(pbReceipt.BlockHash)
	r.BlockNumber = pbReceipt.BlockNumber
	r.TxIndex = pbReceipt.TxIndex
	r.From = common.BytesToAddress(pbReceipt.From)

	if len(pbReceipt.To) > 0 {
		toAddr := common.BytesToAddress(pbReceipt.To)
		r.To = toAddr
	}
	r.Amount = new(big.Int)
	r.Amount.SetBytes(pbReceipt.Amount)

	if len(pbReceipt.ContractAddress) > 0 {
		contractAddr := common.BytesToAddress(pbReceipt.ContractAddress)
		r.ContractAddress = contractAddr
	}
	r.CumulativeGasUsed = pbReceipt.CumulativeGasUsed
	r.GasUsed = pbReceipt.GasUsed
	r.PostState = common.BytesToHash(pbReceipt.PostState)

	r.Logs = make([]*log.Log, len(pbReceipt.Logs))
	for i, pbLog := range pbReceipt.Logs {
		log := &log.Log{}
		log.FromProto(pbLog)
		r.Logs[i] = log
	}

	r.LogsBloom = pbReceipt.LogsBloom
	r.Status = pbReceipt.Status
	return nil
}

func (r *Receipt) Data() *ReceiptData {
	// Convert logs to LogData
	logDataList := make([]*log.LogData, len(r.Logs))
	for i, l := range r.Logs {
		logDataList[i] = &log.LogData{
			Address: l.Address,
			Topics:  l.Topics,
			Data:    l.Data,
		}
	}
	return &ReceiptData{
		Status:            r.Status,
		CumulativeGasUsed: r.CumulativeGasUsed,
		LogsBloom:         r.LogsBloom,
		Logs:              logDataList,
	}
}

func (r *Receipt) Hash() common.Hash {
	// Construct ReceiptData
	receiptData := r.Data()
	// Marshal and hash
	b, _ := proto.Marshal(receiptData.Proto())
	return crypto.Keccak256Hash(b)
}

func (r *Receipt) Marshal() ([]byte, error) {
	return proto.Marshal(r.Proto())
}

func (r *Receipt) String() string {
	return fmt.Sprintf(
		"Receipt{\n  TxHash: %s\n  BlockHash: %s\n  BlockNumber: %d\n  TxIndex: %d\n  From: %s\n  To: %s\n  Amount: %s\n  ContractAddress: %s\n  CumulativeGasUsed: %d\n  GasUsed: %d\n  Status: %d\n  PostState: %s\n  Logs: %v entries\n}",
		r.TxHash.Hex(),
		r.BlockHash.Hex(),
		r.BlockNumber,
		r.TxIndex,
		r.From.Hex(),
		r.To.Hex(),
		r.Amount.String(),
		r.ContractAddress.Hex(),
		r.CumulativeGasUsed,
		r.GasUsed,
		r.Status,
		r.PostState.Hex(),
		r.Logs,
	)
}

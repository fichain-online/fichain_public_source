package log

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"

	"FichainCore/common"
	"FichainCore/crypto"
	pb "FichainCore/proto"
)

type Log struct {
	Address     common.Address // 20 bytes
	Topics      []common.Hash  // 32 bytes each
	Data        []byte
	LogIndex    uint32
	BlockNumber uint64
	TxHash      common.Hash
	TxIndex     uint32
	BlockHash   common.Hash
	Removed     bool
}

// Convert to protobuf
func (l *Log) Proto() *pb.Log {
	topics := make([][]byte, len(l.Topics))
	for i, topic := range l.Topics {
		topics[i] = topic.Bytes()
	}

	return &pb.Log{
		Address:     l.Address.Bytes(),
		Topics:      topics,
		Data:        l.Data,
		LogIndex:    l.LogIndex,
		BlockNumber: l.BlockNumber,
		TxHash:      l.TxHash.Bytes(),
		TxIndex:     l.TxIndex,
		BlockHash:   l.BlockHash.Bytes(),
		Removed:     l.Removed,
	}
}

// Load from protobuf
func (l *Log) FromProto(pbLog *pb.Log) {
	l.Address = common.BytesToAddress(pbLog.Address)

	l.Topics = make([]common.Hash, len(pbLog.Topics))
	for i, topic := range pbLog.Topics {
		l.Topics[i] = common.BytesToHash(topic)
	}

	l.Data = pbLog.Data
	l.LogIndex = pbLog.LogIndex
	l.BlockNumber = pbLog.BlockNumber
	l.TxHash = common.BytesToHash(pbLog.TxHash)
	l.TxIndex = pbLog.TxIndex
	l.BlockHash = common.BytesToHash(pbLog.BlockHash)
	l.Removed = pbLog.Removed
}

func (l *Log) Hash() common.Hash {
	logData := &LogData{
		Address: l.Address,
		Topics:  l.Topics,
		Data:    l.Data,
	}
	b, _ := proto.Marshal(logData.Proto())
	return crypto.Keccak256Hash(b)
}

func (l *Log) String() string {
	topics := make([]string, len(l.Topics))
	for i, topic := range l.Topics {
		topics[i] = topic.Hex()
	}

	return fmt.Sprintf(`Log{
  Address:     %s,
  Topics:      [%s],
  Data:        %x,
  LogIndex:    %d,
  BlockNumber: %d,
  TxHash:      %s,
  TxIndex:     %d,
  BlockHash:   %s,
  Removed:     %t
}`,
		l.Address.Hex(),
		strings.Join(topics, ", "),
		l.Data,
		l.LogIndex,
		l.BlockNumber,
		l.TxHash.Hex(),
		l.TxIndex,
		l.BlockHash.Hex(),
		l.Removed,
	)
}

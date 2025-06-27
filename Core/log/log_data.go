package log

import (
	"FichainCore/common"
	pb "FichainCore/proto"
)

type LogData struct {
	Address common.Address // 20 bytes
	Topics  []common.Hash  // 32 bytes each
	Data    []byte
}

// Convert to protobuf
func (l *LogData) Proto() *pb.LogData {
	topics := make([][]byte, len(l.Topics))
	for i, topic := range l.Topics {
		topics[i] = topic.Bytes()
	}

	return &pb.LogData{
		Address: l.Address.Bytes(),
		Topics:  topics,
		Data:    l.Data,
	}
}

// Load from protobuf
func (l *LogData) FromProto(pbLog *pb.LogData) {
	l.Address = common.BytesToAddress(pbLog.Address)

	l.Topics = make([]common.Hash, len(pbLog.Topics))
	for i, topic := range pbLog.Topics {
		l.Topics[i] = common.BytesToHash(topic)
	}

	l.Data = pbLog.Data
}

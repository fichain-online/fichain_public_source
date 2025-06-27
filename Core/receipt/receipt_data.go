package receipt

import (
	"FichainCore/log"
	pb "FichainCore/proto"
)

type ReceiptData struct {
	Status            uint32
	CumulativeGasUsed uint64
	LogsBloom         []byte // 256 bytes
	Logs              []*log.LogData
}

func (r *ReceiptData) Proto() *pb.ReceiptData {
	logs := make([]*pb.LogData, len(r.Logs))
	for i, log := range r.Logs {
		logs[i] = log.Proto()
	}

	return &pb.ReceiptData{
		Status:            r.Status,
		CumulativeGasUsed: r.CumulativeGasUsed,
		LogsBloom:         r.LogsBloom,
		Logs:              logs,
	}
}

func (r *ReceiptData) FromProto(pbData *pb.ReceiptData) {
	r.Status = pbData.Status
	r.CumulativeGasUsed = pbData.CumulativeGasUsed
	r.LogsBloom = pbData.LogsBloom

	r.Logs = make([]*log.LogData, len(pbData.Logs))
	for i, pbLog := range pbData.Logs {
		log := &log.LogData{}
		log.FromProto(pbLog)
		r.Logs[i] = log
	}
}

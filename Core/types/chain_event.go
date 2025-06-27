package types

import (
	"google.golang.org/protobuf/proto"

	pb "FichainCore/proto"
)

type PbChainEventWrap struct {
	PbEvent *pb.ChainEvent
}

func (b *PbChainEventWrap) Proto() proto.Message {
	return b.PbEvent
}

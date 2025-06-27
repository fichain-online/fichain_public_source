package message

import (
	"google.golang.org/protobuf/proto"

	pb "FichainCore/proto"
)

// Ping represents a Ping message sent by a node
type Ping struct {
	NodeID    string
	Timestamp int64
}

// Proto converts Ping to protobuf message
func (p *Ping) Proto() proto.Message {
	return &pb.Ping{
		NodeId:    p.NodeID,
		Timestamp: p.Timestamp,
	}
}

// FromProto populates Ping from a protobuf Ping message
func (p *Ping) FromProto(pbPing *pb.Ping) error {
	if pbPing == nil {
		return nil
	}
	p.NodeID = pbPing.NodeId
	p.Timestamp = pbPing.Timestamp
	return nil
}

// Pong represents a Pong response from a node
type Pong struct {
	NodeID    string
	Timestamp int64
}

// Proto converts Pong to protobuf message
func (p *Pong) Proto() proto.Message {
	return &pb.Pong{
		NodeId:    p.NodeID,
		Timestamp: p.Timestamp,
	}
}

// FromProto populates Pong from a protobuf Pong message
func (p *Pong) FromProto(pbPong *pb.Pong) error {
	if pbPong == nil {
		return nil
	}
	p.NodeID = pbPong.NodeId
	p.Timestamp = pbPong.Timestamp
	return nil
}

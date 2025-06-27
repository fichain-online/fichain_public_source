package message

import (
	"google.golang.org/protobuf/proto"

	pb "FichainCore/proto"
)

// PeerList represents a list of peer addresses
type PeerList struct {
	Addresses []string
}

// Proto converts PeerList to protobuf PeerList
func (p *PeerList) Proto() proto.Message {
	return &pb.PeerList{
		Addresses: p.Addresses,
	}
}

// FromProto populates PeerList from protobuf PeerList
func (p *PeerList) FromProto(pbPeerList *pb.PeerList) error {
	if pbPeerList == nil {
		return nil
	}
	p.Addresses = pbPeerList.Addresses
	return nil
}

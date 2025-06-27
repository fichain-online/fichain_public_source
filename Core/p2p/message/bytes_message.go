package message

import (
	"errors"

	"google.golang.org/protobuf/proto"

	pb "FichainCore/proto"
)

// HandshakeInit is sent by the client to initiate a handshake
type BytesMessage struct {
	Data []byte
}

// Proto converts HandshakeInit to protobuf format
func (h *BytesMessage) Proto() proto.Message {
	return &pb.BytesMessage{
		Data: h.Data,
	}
}

// FromProto populates HandshakeInit from a protobuf message
func (h *BytesMessage) FromProto(pbMsg *pb.BytesMessage) error {
	if pbMsg == nil {
		return errors.New("missing data")
	}
	h.Data = pbMsg.Data
	return nil
}

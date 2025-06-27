package message

import (
	"google.golang.org/protobuf/proto"

	pb "FichainCore/proto"
)

// HandshakeInit is sent by the client to initiate a handshake
type HandshakeInit struct {
	WalletAddress []byte
	Payload       []byte
}

// Proto converts HandshakeInit to protobuf format
func (h *HandshakeInit) Proto() proto.Message {
	return &pb.HandshakeInit{
		WalletAddress: h.WalletAddress,
		Payload:       h.Payload,
	}
}

// FromProto populates HandshakeInit from a protobuf message
func (h *HandshakeInit) FromProto(pbMsg *pb.HandshakeInit) error {
	if pbMsg == nil {
		return nil
	}
	h.WalletAddress = pbMsg.WalletAddress
	h.Payload = pbMsg.Payload
	return nil
}

// HandshakeAck is sent by the server in response, includes signature
type HandshakeAck struct {
	WalletAddress []byte
	Payload       []byte
	Signature     []byte
}

// Proto converts HandshakeAck to protobuf format
func (h *HandshakeAck) Proto() proto.Message {
	return &pb.HandshakeAck{
		WalletAddress: h.WalletAddress,
		Payload:       h.Payload,
		Signature:     h.Signature,
	}
}

// FromProto populates HandshakeAck from a protobuf message
func (h *HandshakeAck) FromProto(pbMsg *pb.HandshakeAck) error {
	if pbMsg == nil {
		return nil
	}
	h.WalletAddress = pbMsg.WalletAddress
	h.Payload = pbMsg.Payload
	h.Signature = pbMsg.Signature
	return nil
}

// HandshakeConfirm is sent by the client to confirm the handshake
type HandshakeConfirm struct {
	Signature []byte
}

// Proto converts HandshakeConfirm to protobuf format
func (h *HandshakeConfirm) Proto() proto.Message {
	return &pb.HandshakeConfirm{
		Signature: h.Signature,
	}
}

// FromProto populates HandshakeConfirm from a protobuf message
func (h *HandshakeConfirm) FromProto(pbMsg *pb.HandshakeConfirm) error {
	if pbMsg == nil {
		return nil
	}
	h.Signature = pbMsg.Signature
	return nil
}

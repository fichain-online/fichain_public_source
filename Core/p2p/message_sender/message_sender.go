package message_sender

import (
	"errors"
	"time"

	logger "github.com/HendrickPhan/golang-simple-logger"

	"FichainCore/common"
	"FichainCore/config"
	"FichainCore/p2p"
	"FichainCore/p2p/lookup_table"
	"FichainCore/p2p/message"
)

type MessageSender struct {
	lookupTable *lookup_table.LookupTable
}

// NewMessageSender initializes and returns a new MessageSender
func NewMessageSender(lt *lookup_table.LookupTable) *MessageSender {
	return &MessageSender{
		lookupTable: lt,
	}
}

// SendToPeer sends a message directly to a TcpPeer
func (ms *MessageSender) SendToPeer(p p2p.Peer, msg *message.Message) error {
	if p == nil {
		return errors.New("peer is nil")
	}
	return p.Send(msg)
}

// SendToAddress looks up a peer by address and sends the message
func (ms *MessageSender) SendToAddress(addr common.Address, msg *message.Message) error {
	logger.DebugP("Message sender lookup table", ms)
	p, ok := ms.lookupTable.Get(addr)
	if !ok {
		return errors.New("peer not found for address")
	}
	return ms.SendToPeer(p, msg)
}

// SendMessageToPeer builds and sends a structured message directly to a peer
func (ms *MessageSender) SendMessageToPeer(
	p p2p.Peer,
	msgType string,
	payload message.HaveProto,
) error {
	msg := &message.Message{
		Header: &message.Header{
			Version:     1,
			SenderID:    config.GetConfig().NodeID,
			MessageType: msgType,
			Timestamp:   time.Now().Unix(),
			Signature:   []byte{}, // Optionally fill signature later
		},
		Payload: payload,
	}
	logger.DebugP("[MessageSender] sent to peer", p.Address(), msgType)

	return ms.SendToPeer(p, msg)
}

// SendMessageToAddress builds and sends a Message with provided type and payload
func (ms *MessageSender) SendMessageToAddress(
	addr common.Address,
	msgType string,
	payload message.HaveProto,
) error {
	msg := &message.Message{
		Header: &message.Header{
			Version:     1,
			SenderID:    config.GetConfig().NodeID,
			MessageType: msgType,
			Timestamp:   time.Now().Unix(),
			Signature:   []byte{}, // Optionally fill signature later
		},
		Payload: payload,
	}

	return ms.SendToAddress(addr, msg)
}

package message

import (
	"fmt"
	"time"

	logger "github.com/HendrickPhan/golang-simple-logger"
	"google.golang.org/protobuf/proto"

	"FichainCore/call_data"
	pb "FichainCore/proto"
	"FichainCore/receipt"
	"FichainCore/transaction"
	"FichainCore/types"
)

const (
	MessagePing     = "ping"
	MessagePong     = "pong"
	MessagePeerList = "peer_list"

	MessageHandshakeInit    = "handshake_init"
	MessageHandshakeAck     = "handshake_ack"
	MessageHandshakeConfirm = "handshake_confirm"

	MessageSendTransaction = "send_transaction"

	MessageCallSmartContract = "call_smart_contract"
	MessageCallResult        = "call_result"

	//
	MessageGetBalance = "get_balance"
	MessageBalance    = "balance"
	//
	MessageGetNonce = "get_nonce"
	MessageNonce    = "nonce"

	MessageGetReceipt = "get_receipt"
	MessageReceipt    = "receipt"

	MessageGetReceipts = "get_receipts"
	MessageReceipts    = "receipts"

	MessageGetValidators = "get_validators"
	MessageValidators    = "validator"

	MessageGetHeadBlock = "get_head_block"
	MessageHeadBlock    = "head_block"

	MessageGetBlock = "get_block"
	MessageBlock    = "block"

	MessageTxMined = "tx_mined"

	//
	MessageChainEvent = "chain_event"
)

type HaveProto interface {
	Proto() proto.Message
}

// Header represents metadata attached to a P2P message
type Header struct {
	Version     uint32
	SenderID    string
	MessageType string
	Timestamp   int64
	Signature   []byte
}

// Message is the high-level wrapper for P2P messages
type Message struct {
	Header  *Header
	Payload HaveProto
}

// Marshal serializes the Message to protobuf bytes
func (m *Message) Marshal() ([]byte, error) {
	pbMsg, err := m.Proto()
	if err != nil {
		return nil, err
	}
	return proto.Marshal(pbMsg)
}

// Unmarshal deserializes protobuf bytes into a Message
func (m *Message) Unmarshal(data []byte) error {
	var pbMsg pb.Message
	err := proto.Unmarshal(data, &pbMsg)
	if err != nil {
		return err
	}
	return m.FromProto(&pbMsg)
}

// Proto converts to a protobuf Message
func (m *Message) Proto() (*pb.Message, error) {
	if m.Header == nil {
		return nil, fmt.Errorf("header is nil")
	}

	pbHeader := &pb.MessageHeader{
		Version:     m.Header.Version,
		SenderId:    m.Header.SenderID,
		MessageType: m.Header.MessageType,
		Timestamp:   m.Header.Timestamp,
		Signature:   m.Header.Signature,
	}

	var payloadBytes []byte
	var err error

	payloadBytes, err = proto.Marshal(m.Payload.Proto())

	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return &pb.Message{
		Header:  pbHeader,
		Payload: payloadBytes,
	}, nil
}

// FromProto populates the Message from a protobuf message
func (m *Message) FromProto(pbMsg *pb.Message) error {
	if pbMsg.Header == nil {
		return fmt.Errorf("message header is nil")
	}

	m.Header = &Header{
		Version:     pbMsg.Header.Version,
		SenderID:    pbMsg.Header.SenderId,
		MessageType: pbMsg.Header.MessageType,
		Timestamp:   pbMsg.Header.Timestamp,
		Signature:   pbMsg.Header.Signature,
	}

	switch pbMsg.Header.MessageType {
	case MessagePing:
		var p pb.Ping
		if err := proto.Unmarshal(pbMsg.Payload, &p); err != nil {
			return err
		}
		ping := &Ping{}
		if err := ping.FromProto(&p); err != nil {
			return err
		}
		m.Payload = ping
	case MessagePong:
		var p pb.Pong
		if err := proto.Unmarshal(pbMsg.Payload, &p); err != nil {
			return err
		}
		pong := &Pong{}
		if err := pong.FromProto(&p); err != nil {
			return err
		}
		m.Payload = pong
	case MessagePeerList:
		var p pb.PeerList
		if err := proto.Unmarshal(pbMsg.Payload, &p); err != nil {
			return err
		}
		peerList := &PeerList{}
		if err := peerList.FromProto(&p); err != nil {
			return err
		}
		m.Payload = peerList
		// --------------------- Handshake ---------------------
	case MessageHandshakeInit:
		var p pb.HandshakeInit
		if err := proto.Unmarshal(pbMsg.Payload, &p); err != nil {
			return err
		}
		init := &HandshakeInit{}
		if err := init.FromProto(&p); err != nil {
			return err
		}
		m.Payload = init
	case MessageHandshakeAck:
		var p pb.HandshakeAck
		if err := proto.Unmarshal(pbMsg.Payload, &p); err != nil {
			return err
		}
		ack := &HandshakeAck{}
		if err := ack.FromProto(&p); err != nil {
			return err
		}
		m.Payload = ack
	case MessageHandshakeConfirm:
		var p pb.HandshakeConfirm
		if err := proto.Unmarshal(pbMsg.Payload, &p); err != nil {
			return err
		}
		confirm := &HandshakeConfirm{}
		if err := confirm.FromProto(&p); err != nil {
			return err
		}
		m.Payload = confirm
		//
	case MessageTxMined:
		fallthrough
	case MessageSendTransaction:
		var p pb.Transaction
		if err := proto.Unmarshal(pbMsg.Payload, &p); err != nil {
			return err
		}
		tx := &transaction.Transaction{}
		if err := tx.FromProto(&p); err != nil {
			return err
		}
		m.Payload = tx
	case MessageCallSmartContract:
		var p pb.CallSmartContractData
		if err := proto.Unmarshal(pbMsg.Payload, &p); err != nil {
			return err
		}
		tx := &call_data.CallSmartContractData{}
		if err := tx.FromProto(&p); err != nil {
			return err
		}
		m.Payload = tx
	case MessageCallResult:
		var p pb.CallSmartContractResponse
		if err := proto.Unmarshal(pbMsg.Payload, &p); err != nil {
			return err
		}
		tx := &call_data.CallSmartContractResponse{}
		if err := tx.FromProto(&p); err != nil {
			return err
		}
		m.Payload = tx
	case MessageReceipt:
		var p pb.Receipt
		if err := proto.Unmarshal(pbMsg.Payload, &p); err != nil {
			return err
		}
		rc := &receipt.Receipt{}
		if err := rc.FromProto(&p); err != nil {
			return err
		}
		m.Payload = rc
	case MessageChainEvent:
		var p pb.ChainEvent
		if err := proto.Unmarshal(pbMsg.Payload, &p); err != nil {
			return err
		}
		rc := &types.PbChainEventWrap{
			PbEvent: &p,
		}
		m.Payload = rc
	case MessageGetReceipt:
		fallthrough
	case MessageGetReceipts:
		fallthrough
	case MessageGetNonce:
		fallthrough
	case MessageNonce:
		fallthrough
	case MessageBalance:
		var p pb.BytesMessage
		if err := proto.Unmarshal(pbMsg.Payload, &p); err != nil {
			return err
		}
		bMessage := &BytesMessage{}
		if err := bMessage.FromProto(&p); err != nil {
			return err
		}
		m.Payload = bMessage
	case MessageReceipts:
		var p pb.Receipts
		if err := proto.Unmarshal(pbMsg.Payload, &p); err != nil {
			return err
		}
		receipts := &receipt.Receipts{}
		if err := receipts.FromProto(&p); err != nil {
			return err
		}
		m.Payload = receipts
	case MessageGetBalance:
		return nil
	// case "block":
	// 	var b pb.Block
	// 	if err := proto.Unmarshal(pbMsg.Payload, &b); err != nil {
	// 		return err
	// 	}
	// 	blk := &block.Block{}
	// 	blk.FromProto(&b)
	// 	m.Payload = blk
	// case "transaction":
	// 	var tx pb.Transaction
	// 	if err := proto.Unmarshal(pbMsg.Payload, &tx); err != nil {
	// 		return err
	// 	}
	// 	t := &transaction.Transaction{}
	// 	t.FromProto(&tx)
	// 	m.Payload = t
	default:
		logger.Warn("Receive unsupport message", pbMsg.Header.MessageType)
		return fmt.Errorf("unknown message type: %s", pbMsg.Header.MessageType)
	}

	return nil
}

// Optional: Stringer for logging/debugging
func (m *Message) String() string {
	if m.Header == nil {
		return "Invalid Message (nil header)"
	}
	return fmt.Sprintf("Message[%s] from %s at %s",
		m.Header.MessageType,
		m.Header.SenderID,
		time.Unix(m.Header.Timestamp, 0).Format(time.RFC3339),
	)
}

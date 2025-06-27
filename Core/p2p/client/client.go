package client

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"time"

	logger "github.com/HendrickPhan/golang-simple-logger"
	"github.com/google/uuid"

	"FichainCore/common"
	"FichainCore/config"
	"FichainCore/crypto"
	"FichainCore/p2p"
	"FichainCore/p2p/message"
	"FichainCore/p2p/peer"
	"FichainCore/signer"
)

// TCPClient is the implementation of the Client interface that manages outbound TCP connections.
type TCPClient struct {
	timeout time.Duration // Timeout for dialing connections
	signer  *signer.Signer
}

// NewTCPClient creates a new instance of a TCPClient with a given timeout.
func NewTCPClient(
	timeout time.Duration,
	signer *signer.Signer,
) *TCPClient {
	return &TCPClient{
		timeout: timeout,
		signer:  signer,
	}
}

// Dial establishes an outbound TCP connection to a peer using the given address.
// It also performs a handshake to ensure the peer is ready.
func (c *TCPClient) Dial(address string) (p2p.Peer, error) {
	// Attempt to dial the connection with the specified timeout
	conn, err := net.DialTimeout("tcp", address, c.timeout)
	if err != nil {
		return nil, err
	}

	// Perform a handshake (You can implement your custom handshake logic here)
	peer, err := c.handshake(conn)
	if err != nil {
		conn.Close() // Close connection if handshake fails
		return nil, err
	}

	// Return the connected peer
	return peer, nil
}

// handshake performs an initial handshake with the peer and returns the corresponding Peer object.
func (c *TCPClient) handshake(conn net.Conn) (p2p.Peer, error) {
	peerID := conn.RemoteAddr().
		String()
	p := peer.NewTCPPeer(conn, peerID)
	// 1. Send initial handshake with address form signer
	walletAddress, err := c.signer.WalletAddress()
	if err != nil {
		return nil, err
	}
	payload := map[string]any{
		"time_stamp": time.Now().Unix(),
		"uuid":       uuid.New().String(),
	} // add more field if needed later
	bPayload, _ := json.Marshal(payload)

	initMsg := &message.HandshakeInit{
		WalletAddress: walletAddress.Bytes(),
		Payload:       bPayload,
	}
	fmtInitMsg := &message.Message{
		Header: &message.Header{
			Version:     config.GetConfig().Version,
			SenderID:    config.GetConfig().NodeID,
			MessageType: message.MessageHandshakeInit,
			Timestamp:   time.Now().Unix(),
			Signature:   []byte{}, // ignore sign
		},
		Payload: initMsg,
	}
	err = p.Send(fmtInitMsg)
	if err != nil {
		slog.Error(fmt.Sprintf("Error when send Handshake init, %v", err))
	}
	// 2. Read server response (wait for HandshakeAck or timeout)
	timeout := time.After(5 * time.Second)

	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("handshake timeout: did not receive HandshakeAck in time")
		default:
			msg, err := p.ReadMessage()
			if err != nil {
				return nil, fmt.Errorf("failed to read handshake response: %w", err)
			}

			if msg.Header.MessageType != message.MessageHandshakeAck {
				continue // skip unrelated messages
			}

			ack, ok := msg.Payload.(*message.HandshakeAck)
			if !ok {
				continue // malformed payload, ignore
			}

			pub, err := crypto.SigToPub(crypto.Keccak256(bPayload), ack.Signature)
			if err != nil {
				return nil, fmt.Errorf("failed to extrack pub from sign: %w", err)
			}
			addr := crypto.PubkeyToAddress(*pub)
			if addr != common.BytesToAddress(ack.WalletAddress) {
				slog.Warn("Invalid sign in HandshakeConfirm message")
				continue // malformed sign, ignore
			}
			// 3. Send final handshake confirmation
			// sign message and send confirmation
			confirmSign, err := c.signer.SignBytes(ack.Payload)
			if err != nil {
				return nil, fmt.Errorf("failed sign confirm: %w", err)
			}
			confirmMsg := &message.HandshakeConfirm{
				Signature: confirmSign.Bytes(),
			}
			fmtConfirmMsg := &message.Message{
				Header: &message.Header{
					Version:     config.GetConfig().Version,
					SenderID:    config.GetConfig().NodeID,
					MessageType: message.MessageHandshakeConfirm,
					Timestamp:   time.Now().Unix(),
					Signature:   []byte{}, // ignore sign
				},
				Payload: confirmMsg,
			}
			p.Send(fmtConfirmMsg)
			logger.Info("Inited connection with ", peerID)
			return p, nil
		}
	}
}

package peer

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	logger "github.com/HendrickPhan/golang-simple-logger"

	"FichainCore/common"
	"FichainCore/p2p"
	"FichainCore/p2p/message"
)

const maxMessageSize = 128 * 1024 * 1024

type TcpPeer struct {
	id       string
	conn     net.Conn
	addr     string
	writer   *bufio.Writer
	reader   *bufio.Reader
	sendLock sync.Mutex
	alive    bool
	done     chan struct{} // Channel to signal when the peer is done

	walletAddress common.Address
}

// Create a new peer from an accepted or dialed TCP connection
func NewTCPPeer(conn net.Conn, id string) p2p.Peer {
	return &TcpPeer{
		id:     id,
		conn:   conn,
		addr:   conn.RemoteAddr().String(),
		writer: bufio.NewWriter(conn),
		reader: bufio.NewReader(conn),
		alive:  true,
		done:   make(chan struct{}), // Initialize the done channel
	}
}

func (p *TcpPeer) ID() string {
	return p.id
}

func (p *TcpPeer) Address() string {
	return p.addr
}

func (p *TcpPeer) Send(msg *message.Message) error {
	if !p.IsAlive() {
		return errors.New("peer is not alive")
	}

	data, err := msg.Marshal()
	if err != nil {
		return err
	}

	p.sendLock.Lock()
	defer p.sendLock.Unlock()

	// --- UPDATED: Start ---
	// Write message length (8 bytes) followed by message data
	length := uint64(len(data))
	header := make([]byte, 8)
	binary.BigEndian.PutUint64(header, length)
	// --- UPDATED: End ---

	if _, err := p.writer.Write(header); err != nil {
		// On a write error, the peer is likely disconnected.
		p.Close()
		return err
	}

	if _, err := p.writer.Write(data); err != nil {
		p.Close()
		return err
	}

	return p.writer.Flush()
}

func (p *TcpPeer) Close() error {
	p.alive = false
	close(p.done) // Signal the done channel when the peer is closed
	return p.conn.Close()
}

func (p *TcpPeer) IsAlive() bool {
	return p.alive
}

// Done returns a channel that signals when the peer is disconnected
func (p *TcpPeer) Done() <-chan struct{} {
	return p.done
}

// ReadLoop continuously reads incoming messages and routes them
func (p *TcpPeer) ReadLoop(router p2p.Router) {
	defer p.Close() // Ensure the peer is closed when the read loop exits

	for {
		// --- UPDATED: Start ---
		header := make([]byte, 8) // Changed from 4 to 8 bytes
		// Use io.ReadFull to prevent short reads from the network.
		if _, err := io.ReadFull(p.reader, header); err != nil {
			// An error here (like EOF) means the connection is closed.
			break
		}

		// Decode the 8-byte header to a uint64 length.
		length := binary.BigEndian.Uint64(header)

		// Protect against excessively large messages.
		if length > maxMessageSize {
			logger.Error(fmt.Sprintf("Message size %d exceeds limit %d", length, maxMessageSize))
			break
		}

		if length == 0 {
			continue // Handle as keep-alive or ignore
		}

		logger.Warn("Message length", length)
		data := make([]byte, length)
		// Use io.ReadFull to ensure we read the entire message body.
		if _, err := io.ReadFull(p.reader, data); err != nil {
			break
		}
		// --- UPDATED: End ---

		msg := &message.Message{}
		if err := msg.Unmarshal(data); err != nil {
			logger.Error("Failed to unmarshal message:", err)
			continue // Don't disconnect for a single malformed message.
		}

		router.Route(p, msg)
	}
}

func (p *TcpPeer) ReadMessage() (*message.Message, error) {
	// --- UPDATED: Start ---
	header := make([]byte, 8) // Changed from 4 to 8 bytes
	// Use io.ReadFull to prevent short reads.
	if _, err := io.ReadFull(p.reader, header); err != nil {
		p.Close()
		return nil, err
	}

	// Decode the 8-byte header.
	length := binary.BigEndian.Uint64(header)

	// Validate message length.
	if length == 0 {
		return nil, fmt.Errorf("invalid message length: %d", length)
	}
	if length > maxMessageSize {
		err := fmt.Errorf("message size %d exceeds limit %d", length, maxMessageSize)
		p.Close()
		return nil, err
	}

	data := make([]byte, length)
	// Use io.ReadFull to read the entire message.
	if _, err := io.ReadFull(p.reader, data); err != nil {
		p.Close()
		return nil, err
	}
	// --- UPDATED: End ---

	msg := &message.Message{}
	if err := msg.Unmarshal(data); err != nil {
		return nil, err
	}

	return msg, nil
}

func (p *TcpPeer) SetWalletAddress(addr common.Address) {
	p.walletAddress = addr
}

func (p *TcpPeer) WalletAddress() common.Address {
	return p.walletAddress
}

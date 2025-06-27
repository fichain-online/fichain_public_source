package peer // Or your chosen package for peer implementations

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"FichainCore/common"
	"FichainCore/p2p"         // Your P2P interfaces
	"FichainCore/p2p/message" // Your Protobuf message definitions
)

// WebSocketPeer represents a peer connected via WebSocket.
type WebSocketPeer struct {
	id         string
	conn       *websocket.Conn
	sendLock   sync.Mutex
	mu         sync.RWMutex // Protects isClosed
	isClosed   bool
	done       chan struct{}
	closeOnce  sync.Once
	remoteAddr string // Network address of the remote peer

	walletAddress common.Address
}

// NewWebSocketPeer creates a new WebSocketPeer.
// remoteAddr is the network address of the remote peer (e.g., "host:port" from HTTP request or dialed URL).
// isOutbound helps differentiate logs or IDs if needed.
func NewWebSocketPeer(conn *websocket.Conn, remoteAddr string, isOutbound bool) p2p.Peer {
	peerID := ""
	if isOutbound {
		peerID = fmt.Sprintf(
			"ws-out-%s-%s",
			remoteAddr,
			conn.LocalAddr().String(),
		) // More unique ID
	} else {
		peerID = fmt.Sprintf("ws-in-%s-%s", remoteAddr, conn.LocalAddr().String())
	}
	return &WebSocketPeer{
		id:         peerID,
		conn:       conn,
		done:       make(chan struct{}),
		remoteAddr: remoteAddr,
		isClosed:   false,
	}
}

// ID returns the unique identifier of the peer.
func (p *WebSocketPeer) ID() string {
	return p.id
}

// Address returns the network address of the peer.
func (p *WebSocketPeer) Address() string {
	return p.remoteAddr
}

// Send sends a Protobuf marshaled message to the peer.
func (p *WebSocketPeer) Send(msg *message.Message) error {
	p.mu.RLock()
	if p.isClosed {
		p.mu.RUnlock()
		return fmt.Errorf("peer %s is not alive or connection closed", p.id)
	}
	p.mu.RUnlock()

	// Marshal the message using Protobuf
	data, err := msg.Marshal() // Assuming msg.Marshal() is your Protobuf marshaling method
	if err != nil {
		return fmt.Errorf("failed to marshal protobuf message for peer %s: %w", p.id, err)
	}

	p.sendLock.Lock()
	defer p.sendLock.Unlock()

	// Set a write deadline
	deadline := time.Now().Add(10 * time.Second) // Example timeout
	if err := p.conn.SetWriteDeadline(deadline); err != nil {
		// If setting deadline fails, the connection might be bad.
		p.Close() // Attempt to clean up
		return fmt.Errorf("failed to set write deadline for peer %s: %w", p.id, err)
	}

	// Send as a binary message. WebSocket handles framing.
	err = p.conn.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		p.Close() // Close on write error
		return fmt.Errorf("failed to write message to peer %s: %w", p.id, err)
	}
	return nil
}

// Close closes the peer connection.
func (p *WebSocketPeer) Close() error {
	var err error
	p.closeOnce.Do(func() {
		p.mu.Lock()
		if p.isClosed {
			p.mu.Unlock()
			return // Already closed
		}
		p.isClosed = true
		p.mu.Unlock()

		close(p.done) // Signal goroutines relying on Done()

		// Send a WebSocket close message to the peer gracefully
		// Set a deadline for the close message itself.
		_ = p.conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
		closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "peer closing")
		if writeErr := p.conn.WriteMessage(websocket.CloseMessage, closeMsg); writeErr != nil {
			// Log if it's not an error that implies the connection is already gone
			if !websocket.IsCloseError(
				writeErr,
				websocket.CloseAbnormalClosure,
				websocket.CloseGoingAway,
			) &&
				writeErr.Error() != "websocket: close sent" &&
				writeErr.Error() != "use of closed network connection" {
				slog.Debug("Error writing WebSocket close message", "peer", p.id, "error", writeErr)
			}
		}

		// Close the underlying WebSocket connection
		err = p.conn.Close()
		slog.Info("WebSocketPeer connection closed", "peer", p.id, "address", p.remoteAddr)
	})
	return err
}

// IsAlive returns true if the peer connection is considered active.
// It checks the isClosed flag.
func (p *WebSocketPeer) IsAlive() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return !p.isClosed
}

// Done returns a channel that is closed when the peer is disconnected.
func (p *WebSocketPeer) Done() <-chan struct{} {
	return p.done
}

// ReadLoop continuously reads messages from the peer and routes them.
func (p *WebSocketPeer) ReadLoop(router p2p.Router) {
	defer func() {
		slog.Info("Exiting ReadLoop, closing peer", "peer", p.id, "address", p.remoteAddr)
		p.Close() // Ensure Close is called when ReadLoop exits
	}()

	// Configure pong handler for keep-alive
	// Max time to wait for a pong after sending a ping
	pongWait := 60 * time.Second      // Example: must receive a pong within 60s of sending a ping
	pingPeriod := (pongWait * 9) / 10 // Send pings slightly more frequently than pongWait

	// Set initial read deadline
	if err := p.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		slog.Error("Failed to set initial read deadline in ReadLoop", "peer", p.id, "error", err)
		return
	}

	p.conn.SetPongHandler(func(string) error {
		slog.Debug("Pong received", "peer", p.id)
		// Reset read deadline upon receiving a pong
		if err := p.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			slog.Warn("Failed to set read deadline on pong", "peer", p.id, "error", err)
			// Not returning error here to keep connection alive if possible,
			// but this indicates a potential issue with the conn state.
		}
		return nil
	})

	// Goroutine to send pings periodically
	pingTicker := time.NewTicker(pingPeriod)
	defer pingTicker.Stop()

	for {
		select {
		case <-pingTicker.C:
			p.mu.RLock()
			isConnClosed := p.isClosed
			p.mu.RUnlock()
			if isConnClosed {
				return
			}
			// Send a ping message
			p.sendLock.Lock() // Use sendLock to avoid concurrent writes with Send()
			err := p.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err != nil {
				slog.Warn("Failed to set write deadline for ping", "peer", p.id, "error", err)
				p.sendLock.Unlock()
				// If cannot set deadline, connection might be bad.
				// Consider closing, but ReadMessage below will likely catch it.
				return // Exit ReadLoop as connection is likely problematic
			}
			if err := p.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				slog.Warn("Failed to send ping", "peer", p.id, "error", err)
				p.sendLock.Unlock()
				return // Exit ReadLoop, will trigger defer p.Close()
			}
			p.sendLock.Unlock()
		default:
			// Proceed to read message
			// The read deadline is managed by the PongHandler and initial SetReadDeadline.
			// No need to set it in every loop iteration here if pings/pongs are active.
			msg, err := p.ReadMessage() // ReadMessage itself will call p.Close() on critical error
			if err != nil {
				// Error handling for ReadMessage
				if websocket.IsCloseError(
					err,
					websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
				) {
					slog.Info(
						"WebSocket connection closed by peer (normal/going away)",
						"peer",
						p.id,
						"error_details",
						err,
					)
				} else if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
					slog.Warn("WebSocket read error (unexpected close)", "peer", p.id, "error_details", err)
				} else {
					// Check if the error indicates the connection was closed by our side already
					p.mu.RLock()
					stillAlive := !p.isClosed
					p.mu.RUnlock()
					if stillAlive { // Only log as error if we didn't initiate the close
						slog.Error("Error reading message in ReadLoop", "peer", p.id, "error_details", err)
					} else {
						slog.Info("ReadLoop: Error reading message on already closing connection", "peer", p.id, "error_details", err)
					}
				}
				return // Exit loop, defer will handle cleanup.
			}

			if msg != nil {
				// Route the message. The router interface is p2p.Router.
				router.Route(p, msg) // router.Route should not block for long.
			}

			// Check if peer is still alive before continuing loop
			p.mu.RLock()
			isConnClosed := p.isClosed
			p.mu.RUnlock()
			if isConnClosed {
				slog.Info("ReadLoop detected peer has been closed, exiting.", "peer", p.id)
				return
			}
		}
	}
}

// ReadMessage reads a single Protobuf message from the peer.
func (p *WebSocketPeer) ReadMessage() (*message.Message, error) {
	p.mu.RLock()
	if p.isClosed {
		p.mu.RUnlock()
		// Return a specific error or websocket.ErrCloseSent if that's more idiomatic for your setup
		return nil, fmt.Errorf(
			"connection already closed for peer %s when calling ReadMessage",
			p.id,
		)
	}
	p.mu.RUnlock()

	// Read a binary message. WebSocket handles framing and length.
	// The read deadline should be managed by the ReadLoop's ping/pong mechanism.
	// If ReadMessage is called outside ReadLoop, a deadline should be set explicitly before calling.
	messageType, data, err := p.conn.ReadMessage()
	if err != nil {
		// If there's a read error, it's good practice to assume the connection is compromised
		// and initiate a close, especially if it's not a clean WebSocket close error.
		if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			slog.Warn(
				"ReadMessage: Non-clean close error, initiating peer close",
				"peer",
				p.id,
				"error",
				err,
			)
		}
		p.Close() // Ensure peer is marked as closed and resources are cleaned up.
		return nil, err
	}

	if messageType != websocket.BinaryMessage {
		// Could also be TextMessage if server/client agreed on it, but Binary is typical for protobuf
		p.Close() // Close on unexpected message type
		return nil, fmt.Errorf(
			"received non-binary message type: %d from peer %s",
			messageType,
			p.id,
		)
	}

	if len(data) == 0 {
		// Empty message might be permissible depending on protocol, but often indicates an issue or keep-alive.
		// For now, let's treat it as an error if expecting actual data.
		// If your protocol allows empty binary messages, handle accordingly.
		// p.Close() // Optionally close
		return nil, fmt.Errorf("received empty binary message from peer %s", p.id)
	}

	// Unmarshal the data using Protobuf
	msg := &message.Message{}
	if err := msg.Unmarshal(data); err != nil { // Assuming msg.Unmarshal is your Protobuf unmarshaling method
		// Don't necessarily close the connection for a single malformed message,
		// but log it. The caller (e.g., ReadLoop) might decide to disconnect.
		slog.Warn(
			"Failed to unmarshal protobuf message",
			"peer",
			p.id,
			"error",
			err,
			"data_len",
			len(data),
		)
		return nil, fmt.Errorf("failed to unmarshal protobuf message from peer %s: %w", p.id, err)
	}

	return msg, nil
}

// Conn returns the underlying websocket connection.
// Expose this with caution. Useful for specific operations like setting deadlines outside ReadLoop (e.g., handshake).
func (p *WebSocketPeer) Conn() *websocket.Conn {
	return p.conn
}

func (p *WebSocketPeer) SetWalletAddress(addr common.Address) {
	p.walletAddress = addr
}

func (p *WebSocketPeer) WalletAddress() common.Address {
	return p.walletAddress
}

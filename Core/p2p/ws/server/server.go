package server // Or your chosen package name, e.g., "p2p"

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"sync"
	"time"

	logger "github.com/HendrickPhan/golang-simple-logger"
	"github.com/google/uuid" // Still used for unique handshake payload elements
	"github.com/gorilla/websocket"

	"FichainCore/common"
	"FichainCore/config" // Standard log for high-level server messages
	"FichainCore/crypto"
	"FichainCore/p2p" // Your P2P interfaces
	"FichainCore/p2p/lookup_table"
	"FichainCore/p2p/message" // Your Protobuf message definitions
	"FichainCore/p2p/ws/peer" // Where your updated WebSocketPeer (using protobuf) resides
	"FichainCore/signer"
)

// WebSocketServer handles incoming WebSocket connections.
// It now aligns with the p2p.Server interface provided earlier.
type WebSocketServer struct {
	// address string // Address is passed to Listen()
	wsPath      string
	peers       map[string]p2p.Peer // Use p2p.Peer interface type
	peerLock    sync.RWMutex        // Changed to RWMutex for potentially better concurrency
	router      p2p.Router          // Use p2p.Router interface type
	lookupTable *lookup_table.LookupTable
	signer      *signer.Signer
	upgrader    websocket.Upgrader
	httpServer  *http.Server
	mu          sync.Mutex // Protects httpServer and isListening state
	isListening bool
}

// NewWebSocketServer creates a new instance of a WebSocketServer.
// Address is no longer taken here, it's passed to Listen.
func NewWebSocketServer(
	wsPath string,
	router p2p.Router, // Use p2p.Router interface type
	lookupTable *lookup_table.LookupTable,
	signer *signer.Signer,
) *WebSocketServer { // Return concrete type, can be assigned to p2p.Server interface
	return &WebSocketServer{
		wsPath:      wsPath,
		peers:       make(map[string]p2p.Peer),
		router:      router,
		signer:      signer,
		lookupTable: lookupTable,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// TODO: Implement proper origin checking in production.
				// For development, allowing all origins is common.
				slog.Debug(
					"WebSocket CheckOrigin",
					"origin",
					r.Header.Get("Origin"),
					"host",
					r.Host,
				)
				return true
			},
		},
	}
}

// Listen starts the HTTP server and listens for WebSocket upgrade requests on the given address.
// This method conforms to the p2p.Server interface.
func (s *WebSocketServer) Listen(address string) error {
	s.mu.Lock()
	if s.isListening {
		s.mu.Unlock()
		return fmt.Errorf(
			"WebSocket server is already listening on %s%s",
			s.httpServer.Addr,
			s.wsPath,
		)
	}

	mux := http.NewServeMux()
	mux.HandleFunc(s.wsPath, s.handleWebSocketUpgrade)

	s.httpServer = &http.Server{
		Addr:    address,
		Handler: mux,
		// Consider adding timeouts for production robustness:
		// ReadTimeout:    15 * time.Second,
		// WriteTimeout:   15 * time.Second,
		// IdleTimeout:    120 * time.Second,
	}
	s.isListening = true
	s.mu.Unlock()

	log.Printf("WebSocket Server is listening on %s%s...", address, s.wsPath)
	err := s.httpServer.ListenAndServe()

	s.mu.Lock()
	s.isListening = false // Reset listening state when ListenAndServe returns
	s.mu.Unlock()

	if err != nil && err != http.ErrServerClosed {
		log.Printf("WebSocket server ListenAndServe error: %v", err)
		return fmt.Errorf("failed to start WebSocket server: %v", err)
	}
	log.Println("WebSocket Server has stopped listening.")
	return nil // Returns nil if server closed gracefully (http.ErrServerClosed)
}

// handleWebSocketUpgrade attempts to upgrade an HTTP connection to a WebSocket connection.
func (s *WebSocketServer) handleWebSocketUpgrade(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Warn("Failed to upgrade to WebSocket", "remote_addr", r.RemoteAddr, "error", err)
		// http.Error(w, "Could not open websocket connection", http.StatusBadRequest) // Optionally send HTTP error
		return
	}
	slog.Info("WebSocket connection upgrade successful", "remote_addr", r.RemoteAddr)

	// Pass r.RemoteAddr for peer identification and the Address() method
	// The actual peer.ID() will be generated inside NewWebSocketPeer
	go s.handleConnection(conn, r.RemoteAddr)
}

// handleConnection processes an incoming WebSocket connection, performs a handshake, and registers the peer.
func (s *WebSocketServer) handleConnection(wsConn *websocket.Conn, remoteAddr string) {
	// Perform the server-side handshake with the new peer.
	// remoteAddr is the network address (ip:port) from the HTTP request.
	// NewWebSocketPeer will use this for its Address() method and to generate a unique ID.
	newPeer, err := s.performServerHandshake(wsConn, remoteAddr)
	if err != nil {
		slog.Error("WebSocket handshake failed", "remote_addr", remoteAddr, "error", err)
		// Ensure the WebSocket connection is closed if handshake fails.
		// The peer.Close() might not have been called if peer object wasn't fully created.
		wsConn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "handshake failed"))
		wsConn.Close()
		return
	}

	// If handshake is successful, newPeer is a valid p2p.Peer.
	// RegisterPeer will manage its lifecycle (start ReadLoop, handle Done channel).
	s.RegisterPeer(newPeer)
}

// performServerHandshake performs the server-side of the handshake.
// It now uses Protobuf for message payloads.
func (s *WebSocketServer) performServerHandshake(
	wsConn *websocket.Conn,
	remoteAddrForID string, // This is the IP:Port string from HTTP request
) (p2p.Peer, error) {
	// Create a new WebSocketPeer. False indicates it's an inbound connection.
	// The peer.NewWebSocketPeer is the one updated to use Protobuf.
	// It needs the wsConn and the remoteAddr string.
	tempPeer := peer.NewWebSocketPeer(wsConn, remoteAddrForID, false).(interface {
		p2p.Peer
		Conn() *websocket.Conn // Assert that it also has Conn() method for handshake deadlines
	})

	// 1. Wait for initial HandshakeInit message
	// Set a specific read deadline for this handshake message
	// Use tempPeer.Conn() if WebSocketPeer exposes the underlying connection for such specific operations
	if err := tempPeer.Conn().SetReadDeadline(time.Now().Add(15 * time.Second)); err != nil {
		return nil, fmt.Errorf("failed to set read deadline for handshake init: %w", err)
	}
	msg, err := tempPeer.ReadMessage()                                   // ReadMessage now expects Protobuf
	if err := tempPeer.Conn().SetReadDeadline(time.Time{}); err != nil { // Clear deadline
		slog.Warn("Failed to clear read deadline after handshake init", "error", err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to read HandshakeInit from %s: %w", remoteAddrForID, err)
	}

	if msg.Header.MessageType != message.MessageHandshakeInit { // Assuming Protobuf enum name
		return nil, fmt.Errorf(
			"expected HandshakeInit, got %s from %s",
			msg.Header.MessageType,
			remoteAddrForID,
		)
	}

	// The payload should be correctly unmarshaled by WebSocketPeer.ReadMessage into msg.Payload
	initMsg, ok := msg.Payload.(*message.HandshakeInit)
	if !ok || initMsg == nil {
		// This check might be redundant if your Protobuf Unmarshal in Peer guarantees correct type or error.
		// If msg.GetPayload() returns the specific type directly (e.g. *message.HandshakeInit), adjust access.
		return nil, fmt.Errorf(
			"malformed or unexpected HandshakeInit payload from %s",
			remoteAddrForID,
		)
	}
	slog.Info(
		"Received WebSocket HandshakeInit",
		"peer_id_attempt",
		tempPeer.ID(),
		"client_wallet_addr",
		common.BytesToAddress(initMsg.WalletAddress).Hex(),
	)

	// 2. Send HandshakeAck
	// Server signs the payload received in HandshakeInit (initMsg.Payload).
	// initMsg.Payload is already []byte from Protobuf.
	ackSign, err := s.signer.SignBytes(initMsg.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to sign ack for %s: %w", remoteAddrForID, err)
	}
	myServerWalletAddress, err := s.signer.WalletAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to get server wallet address: %w", err)
	}

	// Create the payload for HandshakeAck. This is the challenge for the client.
	// This payload needs to be bytes. If it's structured, marshal it.
	// For consistency with previous JSON example, let's assume it's a simple unique byte array.
	// If this payload itself needs to be structured (e.g. for client to parse), define a proto for it.

	payload := map[string]any{
		"time_stamp": time.Now().Unix(),
		"uuid":       uuid.New().String(),
	} // add more field if needed later

	bPayload, _ := json.Marshal(payload)

	ackMsgPayload := &message.HandshakeAck{
		WalletAddress: myServerWalletAddress.Bytes(),
		Payload:       bPayload,        // This is the challenge client must sign
		Signature:     ackSign.Bytes(), // Server's signature over client's initMsg.Payload
	}

	ackFullMsg := &message.Message{
		Header: &message.Header{
			Version:     config.GetConfig().Version, // Use your config
			SenderID:    config.GetConfig().NodeID,  // Server's NodeID
			MessageType: message.MessageHandshakeAck,
			Timestamp:   time.Now().Unix(),
		},
		Payload: ackMsgPayload, // Set Protobuf oneof
	}

	if err := tempPeer.Send(ackFullMsg); err != nil { // Send now sends Protobuf
		return nil, fmt.Errorf("failed to send HandshakeAck to %s: %w", remoteAddrForID, err)
	}
	slog.Info("Sent WebSocket HandshakeAck", "peer_id_attempt", tempPeer.ID())

	// 3. Wait for HandshakeConfirm
	if err := tempPeer.Conn().SetReadDeadline(time.Now().Add(15 * time.Second)); err != nil {
		return nil, fmt.Errorf("failed to set read deadline for handshake confirm: %w", err)
	}
	msgConfirm, err := tempPeer.ReadMessage()                            // ReadMessage expects Protobuf
	if err := tempPeer.Conn().SetReadDeadline(time.Time{}); err != nil { // Clear deadline
		slog.Warn("Failed to clear read deadline after handshake confirm", "error", err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to read HandshakeConfirm from %s: %w", remoteAddrForID, err)
	}

	if msgConfirm.Header.MessageType != message.MessageHandshakeConfirm {
		return nil, fmt.Errorf(
			"expected HandshakeConfirm, got %s from %s",
			msgConfirm.Header.MessageType,
			remoteAddrForID,
		)
	}

	confirmMsg, ok := msgConfirm.Payload.(*message.HandshakeConfirm) // Protobuf oneof access
	if !ok || confirmMsg == nil {
		return nil, fmt.Errorf(
			"malformed HandshakeConfirm payload from %s, type was %T",
			remoteAddrForID,
			msgConfirm.Payload.Proto().ProtoReflect().Type(),
		)
	}
	slog.Info("Received WebSocket HandshakeConfirm", "peer_id_attempt", tempPeer.ID())

	// Verify client's signature in HandshakeConfirm.
	// It should be over the `challengeBytes` (HandshakeAck.Payload) server sent.
	logger.DebugP("verify hash", hex.EncodeToString(crypto.Keccak256(bPayload)))
	logger.DebugP("received signature", hex.EncodeToString(confirmMsg.Signature))
	pubKey, err := crypto.SigToPub(crypto.Keccak256(bPayload), confirmMsg.Signature)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to recover pubkey from HandshakeConfirm signature from %s: %w",
			remoteAddrForID,
			err,
		)
	}
	clientAuthenticatedWalletAddr := crypto.PubkeyToAddress(*pubKey)
	expectedClientWalletAddr := common.BytesToAddress(
		initMsg.WalletAddress,
	) // From the initial verified message

	if clientAuthenticatedWalletAddr != expectedClientWalletAddr {
		return nil, fmt.Errorf(
			"handshake confirmation signature mismatch for %s: expected %s, got %s from signature",
			remoteAddrForID,
			expectedClientWalletAddr.Hex(),
			clientAuthenticatedWalletAddr.Hex(),
		)
	}

	// Handshake successful
	slog.Info(
		"WebSocket Handshake successful",
		"peer_id",
		tempPeer.ID(),
		"client_wallet_address",
		clientAuthenticatedWalletAddr.Hex(),
	)
	tempPeer.SetWalletAddress(common.BytesToAddress(initMsg.WalletAddress))
	s.lookupTable.Add(clientAuthenticatedWalletAddr, tempPeer) // Add to lookup table

	// At this point, the tempPeer is fully authenticated and ready.
	return tempPeer, nil
}

// Close shuts down the WebSocket server (closes all open connections and the HTTP server).
// This method conforms to the p2p.Server interface.
func (s *WebSocketServer) Close() error {
	s.mu.Lock() // Lock to access httpServer and isListening
	if !s.isListening && s.httpServer == nil {
		s.mu.Unlock()
		slog.Info("WebSocket Server already closed or was never started.")
		return nil
	}
	currentHTTPServer := s.httpServer // Capture current server instance
	s.mu.Unlock()

	slog.Info("WebSocket Server is shutting down...")

	// Close all peer connections
	s.peerLock.Lock()
	peersToClose := make([]p2p.Peer, 0, len(s.peers))
	for _, p := range s.peers {
		peersToClose = append(peersToClose, p)
	}
	s.peers = make(map[string]p2p.Peer) // Clear the map while holding lock
	s.peerLock.Unlock()

	for _, p := range peersToClose {
		slog.Info("Closing connection to peer", "id", p.ID(), "address", p.Address())
		p.Close() // This should trigger the Done() channel monitored by RegisterPeer's goroutine
	}

	// Shutdown the HTTP server
	if currentHTTPServer != nil {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		) // Graceful shutdown timeout
		defer cancel()
		if err := currentHTTPServer.Shutdown(ctx); err != nil {
			slog.Error(
				"WebSocket HTTP server graceful shutdown failed",
				"error",
				err,
				"falling_back_to_close",
				true,
			)
			// Fallback to forceful close if Shutdown fails or times out
			if closeErr := currentHTTPServer.Close(); closeErr != nil {
				slog.Error("WebSocket HTTP server forceful close failed", "error", closeErr)
				return closeErr // Return the forceful close error if all else fails
			}
			return err // Return the shutdown error
		}
		slog.Info("WebSocket HTTP server shut down gracefully.")
	}

	s.mu.Lock()
	s.httpServer = nil // Mark as nil after shutdown
	s.isListening = false
	s.mu.Unlock()

	return nil
}

// RegisterPeer adds an already connected and (usually) authenticated peer to the server's management.
// It starts the peer's ReadLoop and sets up monitoring for its disconnection.
// This method conforms to the p2p.Server interface.
func (s *WebSocketServer) RegisterPeer(p p2p.Peer) {
	if p == nil {
		slog.Warn("Attempted to register a nil peer")
		return
	}
	if !p.IsAlive() {
		slog.Warn("Attempted to register a non-alive peer", "id", p.ID(), "address", p.Address())
		p.Close() // Ensure it's cleaned up if it was somehow created but not alive
		return
	}

	// Ensure it's a WebSocketPeer if this server strictly manages them.
	// This check is more for type safety during development.
	// The p2p.Peer interface should be sufficient.
	if _, ok := p.(interface{ Conn() *websocket.Conn }); !ok { // A loose check for WebSocket peer characteristics
		slog.Warn("Registering a peer that might not be a WebSocketPeer, or doesn't expose Conn()",
			"peer_id", p.ID(), "type", fmt.Sprintf("%T", p))
		// Potentially handle differently or return error if strict typing is required
	}

	s.peerLock.Lock()
	if _, exists := s.peers[p.ID()]; exists {
		s.peerLock.Unlock()
		slog.Warn(
			"Peer already registered, closing new redundant connection",
			"id",
			p.ID(),
			"address",
			p.Address(),
		)
		p.Close() // Close the newly passed-in peer, as one with this ID already exists.
		return
	}
	s.peers[p.ID()] = p
	s.peerLock.Unlock()

	slog.Info("Peer registered with server", "id", p.ID(), "address", p.Address())

	// Start the peer's ReadLoop in a new goroutine.
	// The router will handle incoming messages from this peer.
	go p.ReadLoop(s.router)

	// Start a goroutine to monitor the peer's Done() channel for cleanup.
	go func(peerToRemove p2p.Peer) {
		<-peerToRemove.Done() // Block until the peer's Done channel is closed

		s.peerLock.Lock()
		// Verify that the peer being removed is indeed the one in the map for that ID,
		// in case of race conditions or re-registration attempts.
		if currentPeer, ok := s.peers[peerToRemove.ID()]; ok && currentPeer == peerToRemove {
			delete(s.peers, peerToRemove.ID())
			slog.Info(
				"Peer removed from server's active list",
				"id",
				peerToRemove.ID(),
				"address",
				peerToRemove.Address(),
			)

			// Attempt to remove from lookup table.
			// This requires knowing the key used (client's authenticated wallet address).
			// This logic is highly dependent on how your lookupTable is designed and
			// how/if the authenticated address is retrievable from the p2p.Peer object
			// or if the lookupTable allows removal by Peer object directly.
			// Example: if lookupTable stores WalletAddress->Peer mapping
			// And if peer could store its authenticated wallet address (not standard in p2p.Peer)
			// walletAddr := peerToRemove.GetAuthenticatedWalletAddress() // Fictional method
			// if walletAddr != (common.Address{}) {
			// s.lookupTable.Remove()
			// }
			// OR, if your lookupTable has a RemoveByPeer method:
			// s.lookupTable.RemoveByPeer(peerToRemove)
			// For now, this step is commented out as it needs specific implementation details.
			// slog.Debug("Attempting to remove peer from lookup table", "id", peerToRemove.ID())

		} else if ok {
			// This means a different peer instance is now associated with this ID,
			// or the peer was already removed.
			slog.Warn("Peer for removal was not the current peer in map or already removed", "id", peerToRemove.ID())
		}
		s.peerLock.Unlock()

		log.Printf("Peer disconnected and cleaned up by server: %s", peerToRemove.ID())
	}(p)
}

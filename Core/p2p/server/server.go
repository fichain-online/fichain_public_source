package server

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"

	"FichainCore/common"
	"FichainCore/config"
	"FichainCore/crypto"
	"FichainCore/p2p"
	"FichainCore/p2p/lookup_table"
	"FichainCore/p2p/message"
	"FichainCore/p2p/peer"
	"FichainCore/signer"
)

// TCPServer is the implementation of the Server interface that listens for incoming TCP connections.
type TCPServer struct {
	address     string              // The address to listen on (e.g., "127.0.0.1:8080")
	peers       map[string]p2p.Peer // A map to keep track of connected peers
	peerLock    sync.Mutex          // Mutex to protect access to the peers map
	router      p2p.Router          // Router to route messages
	lookupTable *lookup_table.LookupTable

	signer *signer.Signer
}

// NewTCPServer creates a new instance of a TCPServer to handle incoming connections.
func NewTCPServer(
	address string,
	router p2p.Router,
	lookupTable *lookup_table.LookupTable,
	signer *signer.Signer,
) *TCPServer {
	return &TCPServer{
		address:     address,
		peers:       make(map[string]p2p.Peer),
		router:      router,
		signer:      signer,
		lookupTable: lookupTable,
	}
}

// Listen starts the server and listens for incoming peer connections on the specified address.
func (s *TCPServer) Listen() error {
	// Start listening on the provided address and port
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	defer listener.Close()

	log.Printf("Server is listening on %s...", s.address)

	// Accept incoming connections in a loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		// Handle the connection in a separate goroutine
		go s.handleConnection(conn)
	}
}

// handleConnection processes an incoming TCP connection, performs a handshake, and adds the peer.
func (s *TCPServer) handleConnection(conn net.Conn) {
	// Perform the handshake with the new peer
	peer, err := s.handshake(conn)
	if err != nil {
		log.Printf("Handshake failed: %v", err)
		conn.Close()
		return
	}

	// Add the peer to the server's peer list
	s.peerLock.Lock()
	s.peers[peer.ID()] = peer
	s.peerLock.Unlock()

	// Handle the reading loop from the peer
	go peer.ReadLoop(s.router)

	// Log when the peer connects
	log.Printf("Peer connected: %s", peer.ID())

	// Wait for the peer to disconnect
	<-peer.Done()

	// Once the peer disconnects, clean up
	s.peerLock.Lock()
	delete(s.peers, peer.ID())
	s.peerLock.Unlock()

	log.Printf("Peer disconnected: %s", peer.ID())
}

// handshake performs an initial handshake with the peer (you can add custom logic here).
func (s *TCPServer) handshake(conn net.Conn) (p2p.Peer, error) {
	// Example handshake: Create a new tcpPeer and return it.
	peerID := conn.RemoteAddr().
		String()
	p := peer.NewTCPPeer(conn, peerID)
	// 1. wait for initial handshake

	timeout := time.After(5 * time.Second)
	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("handshake timeout: did not receive HandshakeInit in time")
		default:
			msg, err := p.ReadMessage()
			if err != nil {
				return nil, fmt.Errorf("failed to read handshake response: %w", err)
			}

			if msg.Header.MessageType != message.MessageHandshakeInit {
				continue // skip unrelated messages
			}
			initMsg, ok := msg.Payload.(*message.HandshakeInit)
			if !ok {
				continue // malformed payload, ignore
			}
			slog.Info(fmt.Sprintf("Receive Handshake init from %v", peerID))
			// 2. send ack
			ackSign, err := s.signer.SignBytes(initMsg.Payload)
			if err != nil {
				return nil, fmt.Errorf("failed sign ack: %w", err)
			}
			walletAddress, err := s.signer.WalletAddress()
			if err != nil {
				return nil, fmt.Errorf("failed to get wallet address from signer: %w", err)
			}

			payload := map[string]any{
				"time_stamp": time.Now().Unix(),
				"uuid":       uuid.New().String(),
			} // add more field if needed later

			bPayload, _ := json.Marshal(payload)
			ackMsg := &message.HandshakeAck{
				WalletAddress: walletAddress.Bytes(),
				Payload:       bPayload,
				Signature:     ackSign.Bytes(),
			}
			fmtConfirmMsg := &message.Message{
				Header: &message.Header{
					Version:     config.GetConfig().Version,
					SenderID:    config.GetConfig().NodeID,
					MessageType: message.MessageHandshakeAck,
					Timestamp:   time.Now().Unix(),
					Signature:   []byte{}, // ignore sign
				},
				Payload: ackMsg,
			}
			p.Send(fmtConfirmMsg)
			slog.Info(fmt.Sprintf("Sent Handshake ack to %v", peerID))

			// 3. wait confirm
			timeout = time.After(5 * time.Second)
			for {
				select {
				case <-timeout:
					return nil, fmt.Errorf(
						"handshake timeout: did not receive HandshakeConfirm in time",
					)
				default:
					msg, err := p.ReadMessage()
					if err != nil {
						return nil, fmt.Errorf("failed to read handshake response: %w", err)
					}

					if msg.Header.MessageType != message.MessageHandshakeConfirm {
						continue // skip unrelated messages
					}
					confirmMsg, ok := msg.Payload.(*message.HandshakeConfirm)
					if !ok {
						continue // malformed payload, ignore
					}

					slog.Info(fmt.Sprintf("Receive Handshake confirm from %v", peerID))
					pub, err := crypto.SigToPub(crypto.Keccak256(bPayload), confirmMsg.Signature)
					if err != nil {
						return nil, fmt.Errorf("failed to extrack pub from sign: %w", err)
					}
					addr := crypto.PubkeyToAddress(*pub)
					if addr != common.BytesToAddress(initMsg.WalletAddress) {
						slog.Warn("Invalid sign in HandshakeConfirm message")
						continue // malformed sign, ignore
					}
					// add to lookupTable
					p.SetWalletAddress(common.BytesToAddress(initMsg.WalletAddress))
					s.lookupTable.Add(addr, p)
					return p, nil
				}
			}
		}
	}
}

// Close shuts down the server (closes all open connections).
func (s *TCPServer) Close() error {
	// Iterate over all connected peers and close their connections
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	for _, peer := range s.peers {
		peer.Close()
	}

	log.Println("Server is shutting down.")
	return nil
}

// RegisterPeer adds an already connected peer to the server and starts ReadLoop.
func (s *TCPServer) RegisterPeer(p p2p.Peer) {
	s.peerLock.Lock()
	s.peers[p.ID()] = p
	s.peerLock.Unlock()

	// Start reading loop
	go p.ReadLoop(s.router)

	// Log connection
	log.Printf("Outbound peer registered: %s", p.ID())

	// Wait for disconnection
	go func() {
		<-p.Done()
		s.peerLock.Lock()
		delete(s.peers, p.ID())
		s.peerLock.Unlock()
		log.Printf("Peer disconnected: %s", p.ID())
	}()
}

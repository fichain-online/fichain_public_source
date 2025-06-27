package p2p

import (
	"net"

	"FichainCore/common"
	"FichainCore/p2p/message"
)

// Represents a peer in the network
type Peer interface {
	ID() string
	Address() string
	Send(msg *message.Message) error
	Close() error
	IsAlive() bool
	ReadLoop(router Router)
	ReadMessage() (*message.Message, error)
	Done() <-chan struct{}
	SetWalletAddress(common.Address)
	WalletAddress() common.Address
}

// Defines message routing behavior
type Router interface {
	Route(peer Peer, msg *message.Message)
	RegisterHandler(msgType string, handler func(Peer, *message.Message) error)
}

// Manages outbound connections to other peers
type Client interface {
	Dial(address string) (Peer, error)
}

// Handles incoming TCP connections
type Server interface {
	Listen(address string) error
	Close() error
	RegisterPeer(p Peer)
}

// Manages initial handshake between peers
type Handshaker interface {
	DoHandshake(conn net.Conn) (Peer, error)
}

// Discovers peers to connect to
type PeerDiscovery interface {
	FindPeers() []string
	AddKnownPeer(address string)
}

// Optional: For rate-limiting or spam prevention
type RateLimiter interface {
	Allow(peerID string, msgType string) bool
}

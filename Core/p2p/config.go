package p2p

import (
	"time"
)

// Config holds configuration parameters for the P2P network
type Config struct {
	ListenAddress       string        // Address the server listens on
	MaxPeers            int           // Maximum number of connected peers
	HandshakeTimeout    time.Duration // Timeout for completing a handshake
	MessageReadTimeout  time.Duration // Timeout for reading a message
	MessageWriteTimeout time.Duration // Timeout for sending a message
	EnableRateLimiter   bool          // Enable rate limiting
	Debug               bool          // Enable debug logging
}

// DefaultConfig returns a Config with sane defaults
func DefaultConfig() *Config {
	return &Config{
		ListenAddress:       "0.0.0.0:8080",
		MaxPeers:            100,
		HandshakeTimeout:    5 * time.Second,
		MessageReadTimeout:  10 * time.Second,
		MessageWriteTimeout: 10 * time.Second,
		EnableRateLimiter:   false,
		Debug:               false,
	}
}

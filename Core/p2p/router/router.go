package router

import (
	"sync"

	logger "github.com/HendrickPhan/golang-simple-logger"

	"FichainCore/p2p"
	"FichainCore/p2p/message"
)

// Router implements the Router interface using a map of handlers.
type Router struct {
	handlers map[string]func(p2p.Peer, *message.Message) error
	lock     sync.RWMutex
}

// NewRouter creates a new Router instance.
func NewRouter() *Router {
	return &Router{
		handlers: make(map[string]func(p2p.Peer, *message.Message) error),
	}
}

// Route finds and executes the appropriate handler for the message.
func (r *Router) Route(peer p2p.Peer, msg *message.Message) {
	logger.DebugP("Receive message type", msg.Header.MessageType)
	r.lock.RLock()
	handler, ok := r.handlers[msg.Header.MessageType]
	r.lock.RUnlock()

	if ok {
		err := handler(peer, msg)
		if err != nil {
			logger.Error("[Router] Error when handle message type: ", msg.Header.MessageType, err)
		}
	} else {
		logger.Warn("[Router] No handler registered for message type: ", msg.Header.MessageType)
	}
}

// RegisterHandler registers a handler function for a given message type.
func (r *Router) RegisterHandler(msgType string, handler func(p2p.Peer, *message.Message) error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.handlers[msgType] = handler
}

func (r *Router) RegisterHanlders(handlers map[string]func(p2p.Peer, *message.Message) error) {
	for i, v := range handlers {
		r.RegisterHandler(i, v)
	}
}

package p2p

import (
	"net"
)

// Peer is the object that can
type Peer interface {
	net.Conn
	Send(*Message) error
	StopStreaming()
}

// Transport is anything used to create connection and communication between
// two nodes. It can be TCP, UDP, and websocket
type Transport interface {
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan Message
	Close() error
}

package p2p

import (
	"encoding/gob"
	"net"
	"sync"
)

type TCPPeer struct {
	net.Conn

	// if the peer initiates a dial and retrieve a conn -> outbound == true
	// if the peer accept and retrieve a conn -> outbound == false
	outbound bool

	wg sync.WaitGroup
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
	}
}

func (tp *TCPPeer) Send(msg *Message) error {
	return gob.NewEncoder(tp).Encode(msg)
}

func (tp *TCPPeer) StopStreaming() {
	tp.wg.Done()
}

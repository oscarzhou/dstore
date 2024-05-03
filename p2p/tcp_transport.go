package p2p

import (
	"fmt"
	"log"
	"net"
)

type AddPeerFunc func(peer Peer) error

type TCPTransport struct {
	listenAddress string
	decoder       Decoder
	OnPeer        AddPeerFunc

	listener net.Listener
	msgCh    chan Message
}

func NewTCPTransport(listenAddr string, decoder Decoder) *TCPTransport {
	return &TCPTransport{
		listenAddress: listenAddr,
		msgCh:         make(chan Message),
		decoder:       decoder,
	}
}

func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("failed to dial %s\n", addr)
		return err
	}

	go t.handleConn(conn, true)

	return nil
}

func (t *TCPTransport) Close() error {
	return nil
}

func (t *TCPTransport) Consume() <-chan Message {
	return t.msgCh
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	if t.listener, err = net.Listen("tcp", t.listenAddress); err != nil {
		return err
	}

	go t.startAcceptLoop()
	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			if err == net.ErrClosed {
				log.Printf("elegant exit accept loop")
				return
			}
			log.Println("listener accept error: ", err)
			continue
		}
		log.Printf("[%s] received connection from %s\n", t.listenAddress, conn.RemoteAddr())

		go t.handleConn(conn, false)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	peer := NewTCPPeer(conn, outbound)
	fmt.Printf("new incoming connection %+v\n", peer.RemoteAddr())

	if t.OnPeer != nil {
		if err := t.OnPeer(peer); err != nil {
			log.Printf("on peer failed: %v", err)
			return
		}
	}

	for {
		log.Printf("[%s] is reading message from [%s]\n\n", peer.LocalAddr(), peer.RemoteAddr())
		var msg Message
		if err := t.decoder.Decode(peer, &msg); err != nil {
			log.Printf("1 decode error: %v", err)
			continue
		}

		msg.From = peer.RemoteAddr().String()
		log.Printf("from: %s, read message: %v\n", string(msg.From), msg.Payload)

		peer.wg.Add(1)
		t.msgCh <- msg
		log.Printf("[%s] is waiting\n", string(msg.From))
		peer.wg.Wait()
	}
}

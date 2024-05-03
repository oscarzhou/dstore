package p2p

import (
	"io"
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

func (tp *TCPPeer) Send(r io.Reader) error {
	// msg := &Message{
	// 	Payload: data,
	// }

	// buf := new(bytes.Buffer)

	// if err := gob.NewEncoder(buf).Encode(msg); err != nil {
	// 	return
	// }

	// _, err := tp.Conn.Write(buf.Bytes())
	// if err != nil {
	// 	log.Printf("failed to send data: %v\n", err)
	// 	return
	// }

	// buf := make([]byte, 1024)
	// n, err := r.Read(buf)
	// if err != nil {
	// 	log.Printf("failed to read data before send: %v", err)
	// 	return err
	// }
	// fmt.Printf("send %v", buf[:n])

	// _, err = tp.Write(buf[:n])
	// if err != nil {
	// 	log.Printf("failed to send data: %v", err)
	// 	return err
	// }
	return nil
}

func (tp *TCPPeer) StopStreaming() {
	tp.wg.Done()
}

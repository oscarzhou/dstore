package p2p

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/oscarzhou/dstore/storage"
)

func init() {
	gob.Register(StoreMessage{})
	gob.Register(ReadMessage{})
}

type Server struct {
	Transport

	store   *storage.Store
	peersMu sync.Mutex
	peers   map[string]Peer
	quitCh  chan struct{}

	bootstrapNodes []string
}

func NewServer(transport Transport, store *storage.Store, nodes []string) *Server {
	return &Server{
		Transport:      transport,
		store:          store,
		peers:          make(map[string]Peer),
		quitCh:         make(chan struct{}),
		bootstrapNodes: nodes,
	}
}

func (s *Server) Stop() {
	close(s.quitCh)
}

func (s *Server) GetData(key string) (io.ReadCloser, error) {
	// 1. Get data from local
	if s.store.HasKey(key) {
		return s.store.Get(key)
	}

	// 2. If data doesn't exist in local, get from other peer
	return nil, nil
}

func (s *Server) StoreDataLocal(key string, r io.Reader) error {
	return s.store.Store(key, r)
}

func (s *Server) StoreData(key string, data []byte) error {
	// 1. Store the file to local
	reader := bytes.NewReader(data)
	var buf bytes.Buffer
	r := io.TeeReader(reader, &buf)
	err := s.store.Store(key, r)
	if err != nil {
		return err
	}

	// 2. Broadcast the key and file
	msg := &Message{
		Payload: &StoreMessage{
			Key:      key,
			DataSize: reader.Size(),
		},
	}
	if err := s.broadcast(msg); err != nil {
		return err
	}

	time.Sleep(100 * time.Millisecond)
	return s.stream(bytes.NewReader(buf.Bytes()))
}

func (s *Server) broadcast(msg *Message) error {
	peers := make([]io.Writer, len(s.peers))
	i := 0
	for _, peer := range s.peers {
		peers[i] = peer
		i++
	}
	mw := io.MultiWriter(peers...)

	return gob.NewEncoder(mw).Encode(msg)
}

func (s *Server) stream(r io.Reader) error {
	peers := make([]io.Writer, len(s.peers))
	i := 0
	for _, peer := range s.peers {
		peers[i] = peer
		i++
	}
	mw := io.MultiWriter(peers...)

	_, err := io.Copy(mw, r)
	return err
}

func (s *Server) AddPeer(peer Peer) error {
	s.peersMu.Lock()
	s.peers[peer.RemoteAddr().String()] = peer
	s.peersMu.Unlock()
	return nil
}

func (s *Server) GetPeer(addr string) Peer {
	s.peersMu.Lock()
	defer s.peersMu.Unlock()
	peer, ok := s.peers[addr]
	if !ok {
		return nil
	}
	return peer
}

func (s *Server) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	// join other nodes if there are
	for _, node := range s.bootstrapNodes {
		if err := s.Transport.Dial(node); err != nil {
			log.Printf("failed to dial node from ? to %s", node)
		}
	}

	s.loop()
	return nil
}

func (s *Server) loop() {
	for {
		select {
		case msg := <-s.Transport.Consume():
			log.Printf("loop: from [%s] is receving msg %v\n", msg.From, msg)

			peer := s.GetPeer(msg.From)
			if peer == nil {
				log.Printf("%s peer not found", msg.From)
				continue
			}

			switch v := msg.Payload.(type) {
			case *StoreMessage:
				s.handleStoreMessage(v, peer)
			case *ReadMessage:

			}

		case <-s.quitCh:
			fmt.Println("server quits")
			return
		}
	}
}

func (s *Server) handleStoreMessage(msg *StoreMessage, peer Peer) error {
	log.Printf("start to copy stream, expect size: %d\n", msg.DataSize)
	defer peer.StopStreaming()

	return s.store.Store(string(msg.Key), io.LimitReader(peer, msg.DataSize))
}

func (s *Server) handleReadMessage(msg *ReadMessage, peer Peer) error {

	return nil
}

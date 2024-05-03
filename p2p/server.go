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

func (s *Server) StoreData(key string, data []byte) error {
	// 1. Store the file to local
	reader := bytes.NewReader(data)
	err := s.store.Store(key, reader)
	if err != nil {
		return err
	}

	// 2. Broadcast the key and file
	if err := s.broadcast([]byte(key), reader.Size()); err != nil {
		return err
	}

	time.Sleep(100 * time.Millisecond)
	return s.stream(reader)
}

func (s *Server) broadcast(data []byte, dataSize int64) error {
	peers := make([]io.Writer, len(s.peers))
	i := 0
	for _, peer := range s.peers {
		peers[i] = peer
		i++
	}
	mw := io.MultiWriter(peers...)

	msg := &Message{
		Payload: data,
		Size:    dataSize,
	}
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

			peer, ok := s.peers[msg.From]
			if !ok {
				log.Printf("%s peer not found", msg.From)
				continue
			}

			defer peer.StopStreaming()

			log.Printf("start to copy stream, expect size: %d\n", msg.Size)
			err := s.store.Store(string(msg.Payload), io.LimitReader(peer, msg.Size))
			if err != nil {
				log.Println(err)
				continue
			}

		case <-s.quitCh:
			fmt.Println("server quits")
			return
		}
	}
}

package main

import (
	"log"
	"time"

	"github.com/containerd/log"
	"github.com/oscarzhou/dstore/p2p"
	"github.com/oscarzhou/dstore/storage"
)

func makeServer(storeRoot, advertisePort string, joinPorts ...string) *p2p.Server {
	tt := p2p.NewTCPTransport(advertisePort, &p2p.GobDecoder{})

	store := storage.NewStore(storeRoot, storage.CASTransformKeyFunc)
	s := p2p.NewServer(tt, store, joinPorts)
	tt.OnPeer = s.AddPeer

	return s
}

func main() {

	s1 := makeServer("[netstore:main]", ":40000")
	s2 := makeServer("[netstore:30000]", ":30000", ":40000")

	go func() {
		s1.Start()
	}()

	time.Sleep(1 * time.Second)
	go s2.Start()
	time.Sleep(2 * time.Second)

	// data := []byte("This a distributed storage project")
	// err := s1.StoreData("myprivatekey", data)
	// if err != nil {
	// 	log.Fatalf("store data error: %v", err)
	// }

	key := "myprivatekey"
	reader, err := s1.GetData(key)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	err = s1.StoreDataLocal(key, reader)
	if err != nil {
		log.Fatal(err)
	}

	select {}

}

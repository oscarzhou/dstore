package p2p

import (
	"io"
)

type Message struct {
	From    string
	Payload []byte
	Size    int64
}

type StoreMessage struct {
	Key  string
	Data io.Reader
}

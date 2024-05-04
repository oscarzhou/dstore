package p2p

type Message struct {
	From    string
	Payload interface{}
}

type StoreMessage struct {
	Key      string
	DataSize int64
}

type ReadMessage struct {
	Key string
}

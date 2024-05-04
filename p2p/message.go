package p2p

type MessageType int

const (
	IncomingMessageType MessageType = iota
	StreamMessageType
)

type Message struct {
	From    string
	Type    MessageType
	Payload interface{}
}

type StoreMessage struct {
	Key      string
	DataSize int64
}

type ReadMessage struct {
	Key string
}

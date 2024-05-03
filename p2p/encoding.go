package p2p

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
)

type Decoder interface {
	Decode(io.Reader, *Message) error
}

type GobDecoder struct {
}

func (gd *GobDecoder) Decode(r io.Reader, v *Message) error {
	return gob.NewDecoder(r).Decode(v)
}

type NopDecoder struct{}

func (nd *NopDecoder) Decode(r io.Reader, v *Message) error {
	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	if err != nil {
		log.Printf("read error: %v", err)
		return err
	}

	v.Payload = buf[:n]
	v.Size = int64(n)
	fmt.Printf("write %d bytes\n", n)
	return nil
}

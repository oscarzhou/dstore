package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {

	// go func(t *testing.T) {
	// 	listenAddr := ":30000"
	// 	tt := NewTCPTransport(listenAddr)
	// 	assert.Equal(t, tt.listenAddress, listenAddr)
	// 	time.Sleep(2 * time.Second)
	// 	err := tt.Dial(":40000")
	// 	assert.NoError(t, err, "dial error")

	// }(t)
	listenAddr := ":40000"
	tt := NewTCPTransport(listenAddr, &NopDecoder{})
	assert.Equal(t, tt.listenAddress, listenAddr)
	err := tt.ListenAndAccept()
	assert.NoError(t, err, "listening error")

	select {}
}

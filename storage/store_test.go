package storage

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	s := NewStore("[test]", CASTransformKeyFunc)

	key := "myprivatekey"
	buf := []byte("I have private key")
	reader := bytes.NewReader(buf)
	err := s.Store(key, reader)
	assert.NoError(t, err)
	assert.Equal(t, s.HasKey(key), true)

	r, _, err := s.Get(key)
	assert.NoError(t, err)
	defer r.Close()

	data := make([]byte, 1024)
	n, err := r.Read(data)
	assert.NoError(t, err)
	assert.Equal(t, buf, data[:n])
}

package storage

import (
	"crypto/sha1"
	"encoding/hex"
	"path"
)

func CASTransformKeyFunc(key string) string {
	hashKey := sha1.Sum([]byte(key))
	encodedHashKey := hex.EncodeToString(hashKey[:])
	blockSize := 5
	paths := make([]string, blockSize+1)
	for i := 0; i < blockSize; i++ {
		from, to := i*blockSize, (i+1)*blockSize
		paths[i] = encodedHashKey[from:to]
	}
	paths[blockSize] = encodedHashKey
	tranPath := path.Join(paths...)
	return tranPath
}

package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCASTransformKeyFunc(t *testing.T) {
	expectedPath := "9bbaf/464b0/f6f8a/73536/fc476/9bbaf464b0f6f8a73536fc476bb1bf9587bdab84"
	actualPath := CASTransformKeyFunc("myprivatekey")
	assert.Equal(t, expectedPath, actualPath)
}

package loki

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoki(t *testing.T) {
	l, err := NewLoki("http://localhost:3000")
	assert.NoError(t, err)
	assert.NotNil(t, l)
	assert.Equal(t, "ws://localhost:3000/loki/api/v1/tail", l.Url)
}

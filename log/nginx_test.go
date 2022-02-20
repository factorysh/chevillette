package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNginx(t *testing.T) {

	line := `192.168.1.1 - - [17/Feb/2022:19:25:44 +0100] "GET /uploads/-/system/project/avatar/33/d.png HTTP/2.0" 200 3952 "https://gitlab.example.com/factory/plop" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:98.0) Gecko/20100101 Firefox/98.0"`
	n, err := NewNginxLine("/uploads/")
	assert.NoError(t, err)

	keys, err := n.Log([]byte(line))
	assert.NoError(t, err)
	assert.Equal(t, "192.168.1.1", keys[0])
	assert.Contains(t, keys[1], "Mozilla")
}

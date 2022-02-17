package log

import (
	"testing"

	"github.com/factorysh/chevillette/pattern"
	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	line := `192.168.1.1 - - [17/Feb/2022:19:25:44 +0100] "GET /uploads/-/system/project/avatar/33/d.png HTTP/2.0" 200 3952 "https://gitlab.example.com/factory/plop" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:98.0) Gecko/20100101 Firefox/98.0"`
	p, err := pattern.New(`<ip> - - <_> "<method> <url> <_>" <status> <_> <_> "<_>" <_>`)
	assert.NoError(t, err)
	m := p.Matches([]byte(line))
	ma := make(map[string]string)
	for i, mm := range m {
		ma[p.Names()[i]] = string(mm)
	}
	assert.Equal(t, "192.168.1.1", ma["ip"])
}

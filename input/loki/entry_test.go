package loki

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntry(t *testing.T) {
	raw := []byte(`
	{
		"streams": [
		  {
			"stream": {
				"app": "chevillette"
			},
			"values": [
			  [
				  "1647810544176501185",
				"\"GET /burps HTTP/1.1\" 404 153 \"-\" \"curl/7.82.0-DEV\" \"-\""
			  ]
			]
		  }
		],
		"dropped_entries": [
		  {
			"labels": {
				"beuha": "aussi"
			},
			"timestamp": "1647810544176501185"
		  }
		]
	  }
	
	`)
	var s Tail
	err := json.Unmarshal(raw, &s)
	assert.NoError(t, err)
	assert.Contains(t, s.Streams[0].Entries[0].Line, "/burps")
}

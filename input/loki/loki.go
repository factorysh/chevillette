package loki

import (
	"bytes"
	"context"
	"fmt"
	_url "net/url"
	"time"

	"github.com/davecgh/go-spew/spew"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type Stream struct {
	Stream map[string]string `yaml:"stream"`
	Values [][]string        `yaml:"values"`
}

type DroppedEntry struct {
	Labels    map[string]string `yaml:"labels"`
	Timestamp string            `yaml:"timestamp"`
}

type Tail struct {
	Streams        []Stream       `yaml:"streams"`
	DroppedEntries []DroppedEntry `yaml:"dropped_entries"`
}

type Loki struct {
	Url string
}

func New(url string) (*Loki, error) {
	u, err := _url.Parse(url)
	if err != nil {
		return nil, err
	}
	buff := bytes.Buffer{}
	switch u.Scheme {
	case "http":
		buff.WriteString("ws")
	case "https":
		buff.WriteString("wss")
	default:
		return nil, fmt.Errorf("unknown scheme : %s", u.Scheme)
	}
	buff.WriteString("://")
	buff.WriteString(u.Host)
	if u.Path == "" {
		buff.WriteString("/loki/api/v1/tail")
	} else {
		buff.WriteString(u.Path)
	}
	return &Loki{
		Url: buff.String(),
	}, nil
}

func (l *Loki) Tail(ctx context.Context, query string, delay_for time.Duration, limit int, start time.Time) error {
	c, _, err := websocket.Dial(ctx, "ws://localhost:8080", nil)
	if err != nil {
		return err
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")
	var resp Tail
	for {
		err = wsjson.Read(ctx, c, &resp)
		if err != nil {
			return err
		}
		spew.Dump(resp)
	}
	return nil
}

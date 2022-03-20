package loki

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	_url "net/url"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/websocket"
)

type Stream struct {
	Stream  map[string]string `yaml:"stream"`
	Entries []Entry           `yaml:"values"`
}

type Tail struct {
	Streams        []Stream       `yaml:"streams"`
	DroppedEntries []DroppedEntry `yaml:"dropped_entries"`
}

type Loki struct {
	Url    string
	dialer *websocket.Dialer
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
		Url:    buff.String(),
		dialer: &websocket.Dialer{},
	}, nil
}

func (l *Loki) Tail(ctx context.Context, query string, delayFor time.Duration, limit int, start time.Time) error {
	params := _url.Values{}
	params.Add("query", query)
	params.Add("limit", fmt.Sprintf("%d", limit))
	if delayFor != 0 {
		params.Add("delay_for", fmt.Sprintf("%d", int64(delayFor.Seconds())))
	}
	params.Add("start", fmt.Sprintf("%d", start.UnixNano()))

	u := fmt.Sprintf("%s?%s", l.Url, params.Encode())
	header := &http.Header{}
	c, res, err := l.dialer.DialContext(ctx, u, *header)
	if err != nil {
		buf, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("websocket error : %s (%v)", string(buf), err)
	}

	defer c.Close()
	var resp Tail
	for {
		err = c.ReadJSON(&resp)
		if err != nil {
			return err
		}
		spew.Dump(resp)
	}
	return nil
}

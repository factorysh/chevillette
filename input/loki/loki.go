package loki

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	_url "net/url"
	"time"

	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/websocket"
)

type Stream struct {
	Stream  map[string]string `yaml:"stream"`
	Entries []Entry           `yaml:"values"`
}

type Entry struct {
	Timestamp time.Time
	Line      string
}

type DroppedEntry struct {
	Labels    map[string]string `yaml:"labels"`
	Timestamp time.Time         `yaml:"timestamp"`
}

func (d *DroppedEntry) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp string            `json:"timestamp"`
		Labels    map[string]string `json:"labels"`
	}{
		Timestamp: fmt.Sprintf("%d", d.Timestamp.UnixNano()),
		Labels:    d.Labels,
	})
}

func (d *DroppedEntry) UnmarshalJSON(data []byte) error {
	unmarshal := struct {
		Timestamp string            `json:"timestamp"`
		Labels    map[string]string `json:"labels"`
	}{}
	err := json.Unmarshal(data, &unmarshal)
	if err != nil {
		return err
	}

	t, err := strconv.ParseInt(unmarshal.Timestamp, 10, 64)
	if err != nil {
		return err
	}
	d.Timestamp = time.Unix(0, t)
	d.Labels = unmarshal.Labels
	return nil
}

func (e *Entry) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{
		fmt.Sprintf("%d", e.Timestamp.UnixNano()),
		e.Line,
	})
}

func (e *Entry) UnmarshalJSON(data []byte) error {
	l := make([]string, 0)
	err := json.Unmarshal(data, &l)
	if err != nil {
		return err
	}
	if len(l) != 2 {
		return fmt.Errorf("wrong entry length : %d", len(l))
	}
	t, err := strconv.ParseInt(l[0], 10, 64)
	if err != nil {
		return err
	}
	e.Timestamp = time.Unix(0, t)
	e.Line = l[1]
	return nil
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

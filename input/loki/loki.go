package loki

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	_url "net/url"
	"time"

	"github.com/gorilla/websocket"
)

type Stream struct {
	Stream  map[string]string `json:"stream"`
	Entries []Entry           `json:"values"`
}

type Tail struct {
	Streams        []Stream       `json:"streams"`
	DroppedEntries []DroppedEntry `json:"dropped_entries"`
}

type Loki struct {
	Url      string // websocket url
	lokiRoot string
	dialer   *websocket.Dialer
}

func NewLoki(url string) (*Loki, error) {
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
		Url:      buff.String(),
		lokiRoot: url,
		dialer:   &websocket.Dialer{},
	}, nil
}

func (l *Loki) Tail(ctx context.Context, query string, delayFor time.Duration, limit int, start time.Time, handler func(*Tail) error) error {
	params := _url.Values{}
	params.Add("query", query)
	params.Add("limit", fmt.Sprintf("%d", limit))
	if delayFor != 0 {
		params.Add("delay_for", fmt.Sprintf("%d", int64(delayFor.Seconds())))
	}
	params.Add("start", fmt.Sprintf("%d", start.UnixNano()))

	u := fmt.Sprintf("%s?%s", l.Url, params.Encode())
	header := &http.Header{}

	for {
		ready, err := http.Get(fmt.Sprintf("%s/ready", l.lokiRoot))
		if err != nil {
			log.Println(err)
			time.Sleep(10 * time.Second)
			continue
		}
		if ready.StatusCode != 200 {
			body, err := ioutil.ReadAll(ready.Body)
			if err != nil {
				return err
			}
			err = ready.Body.Close()
			if err != nil {
				return err
			}
			log.Println("Loki is not ready :", string(body))
			time.Sleep(10 * time.Second)
			continue
		}

		c, res, err := l.dialer.DialContext(ctx, u, *header)
		if err != nil {
			buf, _ := ioutil.ReadAll(res.Body)
			log.Printf("Loki websocket error : %s (%v)", string(buf), err)
			time.Sleep(10 * time.Second)
			continue
		}

		defer c.Close()
		var resp Tail
		for {
			err = c.ReadJSON(&resp)
			if err != nil {
				return err
			}
			err = handler(&resp)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

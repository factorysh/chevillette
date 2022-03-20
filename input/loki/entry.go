package loki

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Entry struct {
	Timestamp time.Time
	Line      string
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

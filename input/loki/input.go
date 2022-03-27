package loki

import (
	"context"
	"time"

	"github.com/factorysh/chevillette/log"
	"github.com/factorysh/chevillette/memory"
)

type LokiInput struct {
	Url    string
	Query  string
	line   log.LineReader
	loki   *Loki
	memory *memory.Memory
}

func New(url string, query string, line log.LineReader, memory *memory.Memory) (*LokiInput, error) {
	l, err := NewLoki(url)
	if err != nil {
		return nil, err
	}
	return &LokiInput{
		Url:    url,
		Query:  query,
		line:   line,
		loki:   l,
		memory: memory,
	}, nil
}

func (l *LokiInput) Serve(ctx context.Context) error {
	return l.loki.Tail(ctx, l.Query, 0, 500, time.Now(),
		func(data *Tail) error {
			for _, stream := range data.Streams {
				for _, entry := range stream.Entries {
					keys, err := l.line([]byte(entry.Line))
					if err != nil {
						return err
					}
					l.memory.Set(keys, entry.Timestamp)

				}
			}
			return nil
		})

}

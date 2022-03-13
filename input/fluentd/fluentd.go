package fluentd

import (
	"fmt"
	"time"

	"github.com/factorysh/chevillette/log"
	"github.com/factorysh/chevillette/memory"
	"github.com/factorysh/fluent-server/server"
)

type FluentdInput struct {
	server *server.Server
	tag    string
	line   log.LineReader
	logKey string
	memory memory.Memory
}

func New(tag string, line log.LineReader, memory *memory.Memory) (*FluentdInput, error) {
	f := &FluentdInput{
		tag:    tag,
		line:   line,
		logKey: "log",
		memory: *memory,
	}
	s, err := server.New(func(tag string, ts *time.Time, record map[string]interface{}) error {
		if tag == f.tag {
			fmt.Println(tag, ts, record)
			keys, err := f.line([]byte(record[f.logKey].(string)))
			if err != nil {
				fmt.Println("error", err)
				return nil
			}
			memory.Set(keys, *ts)
			fmt.Println("keys", keys)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	f.server = s
	return f, nil
}

func (f *FluentdInput) ListenAndServe(listen string) error {
	return f.server.ListenAndServe(listen)
}

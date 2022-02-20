package fluentd

import (
	"fmt"
	"time"

	"github.com/factorysh/chevillette/log"
	"github.com/factorysh/fluent-server/server"
)

type FluentdInput struct {
	server *server.Server
	tag    string
	line   log.LineReader
	logKey string
}

func New(tag string, line log.LineReader) *FluentdInput {
	f := &FluentdInput{
		tag:    tag,
		line:   line,
		logKey: "log",
	}
	s := server.New(func(tag string, ts *time.Time, record map[string]interface{}) error {
		if tag == f.tag {
			fmt.Println(tag, ts, record)
			keys, err := f.line([]byte(record[f.logKey].(string)))
			if err != nil {
				fmt.Println("error", err)
				return nil
			}
			fmt.Println("keys", keys)
		}
		return nil
	})
	f.server = s
	return f
}

func (f *FluentdInput) ListenAndServe(listen string) error {
	return f.server.ListenAndServe(listen)
}

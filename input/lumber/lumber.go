package lumber

import (
	_log "log"
	"net"
	"time"

	"github.com/elastic/go-lumber/server"
	"github.com/factorysh/chevillette/conf"
	"github.com/factorysh/chevillette/log"
	"github.com/factorysh/chevillette/memory"
)

type LumberInput struct {
	server server.Server
	line   log.LineReader
	memory *memory.Memory
	stop   chan interface{}
}

func New(line log.LineReader, mem *memory.Memory, cfg conf.Lumber) (*LumberInput, error) {
	l, err := net.Listen("tcp", cfg.Listen)
	if err != nil {
		return nil, err
	}
	s, err := server.NewWithListener(l,
		server.Keepalive(3*time.Minute),
		server.Timeout(30*time.Second),
		server.V2(true),
		server.V1(false),
	)
	if err != nil {
		return nil, err
	}
	return &LumberInput{
		server: s,
		line:   line,
		memory: mem,
		stop:   make(chan interface{}),
	}, nil
}

func (l *LumberInput) Serve() error {
	receive := l.server.ReceiveChan()
	for {
		select {
		case <-l.stop:
			return nil
		case batch := <-receive:
			for _, event := range batch.Events {
				evt, ok := event.(map[string]interface{})
				if !ok {
					continue
				}
				message, ok := evt["message"].(string)
				if !ok {
					continue
				}
				keys, err := l.line([]byte(message))
				if err != nil {
					_log.Println(err)
				}
				// FIXME parce @timestamp
				l.memory.Set(keys, time.Now())
				_log.Println("keys", keys)
			}
			batch.ACK()
		}
	}
}

func (l *LumberInput) Close() error {
	l.stop <- true
	return l.server.Close()
}

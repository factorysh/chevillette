package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/factorysh/chevillette/auth/authrequest"
	"github.com/factorysh/chevillette/conf"
	"github.com/factorysh/chevillette/input/fluentd"
	"github.com/factorysh/chevillette/input/loki"
	"github.com/factorysh/chevillette/log"
	"github.com/factorysh/chevillette/memory"
	"gopkg.in/yaml.v3"
)

func main() {
	cfgPath := os.Getenv("CONFIG")
	if cfgPath == "" {
		cfgPath = "/etc/chevillette.yml"
	}
	c, err := os.Open(cfgPath)
	if err != nil {
		panic(err)
	}
	var cfg conf.Conf
	err = yaml.NewDecoder(c).Decode(&cfg)
	if err != nil {
		panic(err)
	}
	cfg.Default()

	l, err := log.NewNginxLine(cfg.Source.PrefixLines...)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m := memory.New(ctx, cfg.TTL)

	ar := authrequest.New(m)

	if cfg.Fluentd != nil && cfg.Fluentd.Listen != "" {
		f, err := fluentd.New("nginx", l.Log, m, cfg.Fluentd)
		if err != nil {
			panic(err)
		}
		go f.ListenAndServe(cfg.Fluentd.Listen)
	}

	if cfg.Loki != nil {
		lok, err := loki.NewLoki(cfg.Loki.Url)
		if err != nil {
			panic(err)
		}
		go func() {
			err = lok.Tail(context.TODO(), cfg.Loki.Query, 2*time.Second, 1000, time.Now(),
				func(data *loki.Tail) error {
					for _, stream := range data.Streams {
						for _, entry := range stream.Entries {
							slugs, err := l.Log([]byte(entry.Line))
							if err != nil {
								panic(err)
							}
							fmt.Println(slugs)
						}
					}
					return nil
				})
			if err != nil {
				panic(err)
			}
		}()
	}

	err = ar.ListenAndServe(cfg.AuthListen)
	if err != nil {
		panic(err)
	}
}
